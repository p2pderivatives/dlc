package wallet

import (
	"errors"
	"fmt"
	"testing"

	"github.com/btcsuite/btcwallet/walletdb"
	"github.com/btcsuite/btcwallet/wtxmgr"
	"github.com/stretchr/testify/assert"
)

var (
	// Spends: bogus
	// Outputs: 10 BTC
	fakeTxRecordA *wtxmgr.TxRecord

	// Spends: A:0
	// Outputs: 5 BTC, 5 BTC
	fakeTxRecordB *wtxmgr.TxRecord

	exampleBlock100 = makeBlockMeta(100)
)

// Test setup?
// 		create wallet
// 		mine regtest coins?
//  	ListUnspent() to check if we can see the mined coins?

// TestListUnspent() will also need to check different types of scripts
func TestListUnspent(t *testing.T) {
	tearDownFunc, w, db := setupWallet(t)
	defer tearDownFunc()

	_ = fakeUtxos(w, db)

	err := walletdb.View(db, func(tx walletdb.ReadTx) error {
		wtxmgrBucket := tx.ReadBucket(wtxmgrNamespaceKey)
		if wtxmgrBucket == nil {
			return errors.New("missing transaction manager namespace")
		}
		utxos, e := w.ListUnspent()

		assert.Nil(t, e)
		assert.NotNil(t, utxos)
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
	assert.Nil(t, err)

}
