package wallet

import (
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcwallet/waddrmgr"
	"github.com/btcsuite/btcwallet/walletdb"
	"github.com/btcsuite/btcwallet/wtxmgr"
)

// Utxo is a unspend transaction output
type Utxo = btcjson.ListUnspentResult

// ListUnspent returns unspent transactions.
// It asks the running bitcoind instance to return all known utxos for addresses it knows about
func (w *wallet) ListUnspent() (utxos []Utxo, err error) {
	return w.rpc.ListUnspent()
}

// ListUnspent2 also returns unspent transactions except it returns utxos from its own db.
// TODO: add filter
//   Only utxos with address contained the param addresses will be considered.
//   If param addresses is empty, all addresses are considered and there is no
//   filter
func (w *wallet) ListUnspent2() (utxos []Utxo, err error) {
	var results []btcjson.ListUnspentResult
	err = walletdb.View(w.db, func(tx walletdb.ReadTx) error {
		addrmgrNs := tx.ReadBucket(waddrmgrNamespaceKey)
		txmgrNs := tx.ReadBucket(wtxmgrNamespaceKey)

		// filter := len(addresses) != 0

		unspent, e := w.txStore.UnspentOutputs(txmgrNs)
		if e != nil {
			return e
		}

		// utxos = make([]*btcjson.ListUnspentResult, 0, len(unspent))
		for i := range unspent {
			output := unspent[i]
			result := w.credit2ListUnspentResult(output, addrmgrNs)
			// TODO: result might return nil... catch that nil?
			results = append(results, *result)
		}
		return nil
	})
	utxos = results
	return utxos, err
}

func (w *wallet) credit2ListUnspentResult(
	c wtxmgr.Credit,
	addrmgrNs walletdb.ReadBucket) *btcjson.ListUnspentResult {

	syncBlock := w.manager.SyncedTo()

	// TODO: add minconf, maxconf params
	confs := confirms(c.Height, syncBlock.Height)
	// // Outputs with fewer confirmations than the minimum or more
	// // confs than the maximum are excluded.
	// confs := confirms(output.Height, syncBlock.Height)
	// if confs < minconf || confs > maxconf {
	// 	continue
	// }

	// Only mature coinbase outputs are included.
	if c.FromCoinBase {
		target := int32(w.params.CoinbaseMaturity) // make param
		if !confirmed(target, c.Height, syncBlock.Height) {
			// continue
			return nil // maybe?

		}
	}

	acctName := accountName

	result := &btcjson.ListUnspentResult{
		TxID:          c.OutPoint.Hash.String(),
		Vout:          c.OutPoint.Index,
		Account:       acctName,
		ScriptPubKey:  hex.EncodeToString(c.PkScript),
		Amount:        c.Amount.ToBTC(),
		Confirmations: int64(confs),
		Spendable:     true,
	}

	return result
}

// isSpendable determines if given ScriptClass is spendable or not.
// Does NOT support watch-only addresses. This func will need to be rewritten
// to support watch-only addresses
func (w *wallet) isSpendable(sc txscript.ScriptClass, addrs []btcutil.Address,
	addrmgrNs walletdb.ReadBucket) (spendable bool) {
	// At the moment watch-only addresses are not supported, so all
	// recorded outputs that are not multisig are "spendable".
	// Multisig outputs are only "spendable" if all keys are
	// controlled by this wallet.
	//
	// TODO: Each case will need updates when watch-only addrs
	// is added.  For P2PK, P2PKH, and P2SH, the address must be
	// looked up and not be watching-only.  For multisig, all
	// pubkeys must belong to the manager with the associated
	// private key (currently it only checks whether the pubkey
	// exists, since the private key is required at the moment).
scSwitch:
	switch sc {
	case txscript.PubKeyHashTy:
		spendable = true
	case txscript.PubKeyTy:
		spendable = true
	case txscript.WitnessV0ScriptHashTy:
		spendable = true
	case txscript.WitnessV0PubKeyHashTy:
		spendable = true
	case txscript.MultiSigTy:
		for _, a := range addrs {
			_, err := w.manager.Address(addrmgrNs, a)
			if err == nil {
				continue
			}
			if waddrmgr.IsError(err, waddrmgr.ErrAddressNotFound) {
				break scSwitch
			}
			// return err TODO: figure out what to replace the return error
		}
		spendable = true
	}

	return spendable
}

// confirms returns the number of confirmations for a transaction in a block at
// height txHeight (or -1 for an unconfirmed tx) given the chain height
// curHeight.
func confirms(txHeight, curHeight int32) int32 {
	switch {
	case txHeight == -1, txHeight > curHeight:
		return 0
	default:
		return curHeight - txHeight + 1
	}
}

// confirmed checks whether a transaction at height txHeight has met minconf
// confirmations for a blockchain at height curHeight.
func confirmed(minconf, txHeight, curHeight int32) bool {
	return confirms(txHeight, curHeight) >= minconf
}

// SelectUnspent is an implementation of Wallet.SelectUnspent
func (w *wallet) SelectUnspent(
	amt, feePerTxIn, feePerTxOut btcutil.Amount,
) (utxos []Utxo, change btcutil.Amount, err error) {
	var utxosAll []Utxo
	utxosAll, err = w.ListUnspent()
	if err != nil {
		return
	}

	var total btcutil.Amount
	var fee btcutil.Amount
	var utxoAmt btcutil.Amount
	for _, utxo := range utxosAll {
		utxoAmt, err = btcutil.NewAmount(utxo.Amount)
		if err != nil {
			return
		}
		total += utxoAmt
		fee += feePerTxIn
		utxos = append(utxos, utxo)
		if amt+fee == total {
			return
		} else if amt+fee < total {
			change = total - (amt + fee)
			fee += feePerTxOut
			if amt+fee <= total {
				return
			}
		}
	}

	err = fmt.Errorf("Not enough utxos")
	return utxos, change, err
}

// UtxosToTxIns converts utxos to txins
func UtxosToTxIns(utxos []Utxo) ([]*wire.TxIn, error) {
	var txins []*wire.TxIn
	for _, utxo := range utxos {
		txid, err := chainhash.NewHashFromStr(utxo.TxID)
		if err != nil {
			return txins, err
		}
		op := wire.NewOutPoint(txid, utxo.Vout)
		txins = append(txins, wire.NewTxIn(op, nil, nil))
	}
	return txins, nil
}
