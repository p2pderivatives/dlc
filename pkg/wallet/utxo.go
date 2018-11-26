package wallet

import (
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
)

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
