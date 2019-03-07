package wallet

import (
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcwallet/walletdb"
	"github.com/p2pderivatives/dlc/pkg/wallet"
)

const (
	unspentMinConf = 1
	unspentMaxConf = 9999999
)

// ListUnspent returns unspent transactions.
// It asks the running bitcoind instance to return all known utxos for addresses managed by the wallet
func (w *Wallet) ListUnspent() (utxos []wallet.Utxo, err error) {
	var addrs []btcutil.Address
	err = walletdb.View(w.db, func(tx walletdb.ReadTx) error {
		ns := tx.ReadBucket(waddrmgrNamespaceKey)
		appendAddr := func(addr btcutil.Address) error {
			addrs = append(addrs, addr)
			return nil
		}
		return w.manager.ForEachActiveAddress(ns, appendAddr)
	})
	if err != nil {
		return utxos, err
	}

	if addrs == nil {
		return utxos, err
	}

	return w.rpc.ListUnspentMinMaxAddresses(
		unspentMinConf, unspentMaxConf, addrs)
}

// SelectUnspent is an implementation of Wallet.SelectUnspent
func (w *Wallet) SelectUnspent(
	amt, feePerTxIn, feePerTxOut btcutil.Amount,
) (utxos []wallet.Utxo, change btcutil.Amount, err error) {
	var utxosAll []wallet.Utxo
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

// UtxoByTxIn finds utxo by txin
func (w *Wallet) UtxoByTxIn(txin *wire.TxIn) (wallet.Utxo, error) {
	txid := txin.PreviousOutPoint.Hash.String()
	vout := txin.PreviousOutPoint.Index

	utxos, err := w.ListUnspent()
	if err != nil {
		return wallet.Utxo{}, err
	}

	for _, utxo := range utxos {
		if utxo.TxID == txid && utxo.Vout == vout {
			return utxo, nil
		}
	}

	return wallet.Utxo{}, errors.New("utxo isn't found")
}
