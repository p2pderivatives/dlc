package wallet

import (
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcwallet/walletdb"
)

// Utxo is a unspend transaction output
type Utxo = btcjson.ListUnspentResult

// ListUnspent returns unspent transactions.
// TODO: add filter
//   Only utxos with address contained the param addresses will be considered.
//   If param addresses is empty, all addresses are considered and there is no
//   filter
func (w *wallet) ListUnspent() (utxos []*Utxo, err error) {
	var results []*btcjson.ListUnspentResult
	err = walletdb.View(w.db, func(tx walletdb.ReadTx) error {
		addrmgrNs := tx.ReadBucket(waddrmgrNamespaceKey)
		txmgrNs := tx.ReadBucket(wtxmgrNamespaceKey)

		syncBlock := w.manager.SyncedTo()
		// filter := len(addresses) != 0

		unspent, e := w.txStore.UnspentOutputs(txmgrNs)
		if e != nil {
			return e
		}

		// utxos = make([]*btcjson.ListUnspentResult, 0, len(unspent))
		for i := range unspent {
			output := unspent[i]
			result := w.credit2ListUnspentResult(output, syncBlock, addrmgrNs)
			// TODO: result might return nil... catch that nil?
			results = append(results, result)
		}
		return nil
	})
	utxos = results
	return utxos, err
}
