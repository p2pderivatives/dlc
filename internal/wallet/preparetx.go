package wallet

import (
	"fmt"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	_ "github.com/btcsuite/btcwallet/walletdb/bdb" // blank import for bolt db driver
	"github.com/dgarage/dlc-demo/src/dlc"
)

// FundTx adds inputs to a transaction until amount.
func (w *Wallet) FundTx(tx *wire.MsgTx, amount, efee int64) error {
	list, err := w.rpc.ListUnspent()
	if err != nil {
		return err
	}
	outs := []*wire.OutPoint{}
	total := int64(0)
	addfee := int64(0)
	for _, utxo := range list {
		txid, _ := chainhash.NewHashFromStr(utxo.TxID)
		outs = append(outs, wire.NewOutPoint(txid, utxo.Vout))
		a, _ := btcutil.NewAmount(utxo.Amount)
		total += int64(a)
		addfee = int64(len(outs)) * dlc.DlcTxInSize * efee
		if amount+addfee <= total {
			if amount+addfee == total {
				break
			}
			addfee += dlc.DlcTxOutSize * efee
			if amount+addfee <= total {
				break
			}
		}
	}
	if amount+addfee > total {
		return fmt.Errorf("short of bitcoin")
	}
	for _, out := range outs {
		tx.AddTxIn(wire.NewTxIn(out, nil, nil))
	}
	if amount+addfee == total {
		return nil
	}
	change := total - (amount + addfee)
	pkScript := w.P2WPKHpkScript(w.GetFakePublicKey())
	tx.AddTxOut(wire.NewTxOut(change, pkScript))
	return nil
}
