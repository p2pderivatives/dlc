package wallet

import (
	"github.com/btcsuite/btcd/btcjson"
)

// Utxo is a unspend transaction output
type Utxo btcjson.ListUnspentResult

// ListUnspent returns unspent transactions
func (w *wallet) ListUnspent() (utxos []Utxo, err error) {
	return
}
