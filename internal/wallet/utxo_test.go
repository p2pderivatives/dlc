package wallet

import (
	"encoding/binary"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
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
	tearDownFunc, wallet := setupWallet(t)
	defer tearDownFunc()

	utxos := fakeUtxos(wallet)

	syncBlock := wallet.manager.SyncedTo()

	err := walletdb.View(wallet.db, func(tx walletdb.ReadTx) error {
		wtxmgrBucket := tx.ReadBucket(wtxmgrNamespaceKey)
		if wtxmgrBucket == nil {
			return errors.New("missing transaction manager namespace")
		}
		result := wallet.credit2ListUnspentResult(utxos[0], syncBlock, wtxmgrBucket)

		assert.NotNil(t, result)
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
	assert.Nil(t, err)

}

// fakeUtxos creates fake transactions, and inserts them into the provided wallet's db
func fakeUtxos(w *wallet) []wtxmgr.Credit {
	tx := spendOutput(&chainhash.Hash{}, 0, 10e8)
	rec, err := wtxmgr.NewTxRecordFromMsgTx(tx, timeNow())
	if err != nil {
		panic(err)
	}
	fakeTxRecordA = rec

	tx = spendOutput(&fakeTxRecordA.Hash, 0, 5e8, 5e8)
	rec, err = wtxmgr.NewTxRecordFromMsgTx(tx, timeNow())
	if err != nil {
		panic(err)
	}
	fakeTxRecordB = rec

	var utxos []wtxmgr.Credit
	err = walletdb.Update(w.db, func(tx walletdb.ReadWriteTx) error {
		wtxmgrBucket := tx.ReadWriteBucket(wtxmgrNamespaceKey)
		if wtxmgrBucket == nil {
			return errors.New("missing transaction manager namespace")
		}

		_ = w.txStore.InsertTx(wtxmgrBucket, fakeTxRecordA, nil)
		_ = w.txStore.AddCredit(wtxmgrBucket, fakeTxRecordA, nil, 0, false)
		fmt.Printf("created fake credit A\n")

		// Insert a second transaction which spends the output, and creates two
		// outputs.  Mark the second one (5 BTC) as wallet change.
		_ = w.txStore.InsertTx(wtxmgrBucket, fakeTxRecordB, nil)
		_ = w.txStore.AddCredit(wtxmgrBucket, fakeTxRecordB, nil, 1, true)
		fmt.Printf("created fake credit B\n")

		// // Mine each transaction in a block at height 100.
		_ = w.txStore.InsertTx(wtxmgrBucket, fakeTxRecordA, &exampleBlock100)
		_ = w.txStore.InsertTx(wtxmgrBucket, fakeTxRecordB, &exampleBlock100)
		fmt.Printf("mined each transaction\n")

		// Print the one confirmation balance.
		bal, e := w.txStore.Balance(wtxmgrBucket, 1, 100)
		if e != nil {
			fmt.Println(e)
			return nil
		}
		fmt.Println(bal)

		// Fetch unspent outputs.
		utxos, e = w.txStore.UnspentOutputs(wtxmgrBucket)
		if e != nil {
			fmt.Println(e)
		}
		return e
	})
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return utxos
}

func spendOutput(txHash *chainhash.Hash, index uint32, outputValues ...int64) *wire.MsgTx {
	tx := wire.MsgTx{
		TxIn: []*wire.TxIn{
			{
				PreviousOutPoint: wire.OutPoint{Hash: *txHash, Index: index},
			},
		},
	}
	for _, val := range outputValues {
		tx.TxOut = append(tx.TxOut, &wire.TxOut{Value: val})
	}
	return &tx
}

func makeBlockMeta(height int32) wtxmgr.BlockMeta {
	if height == -1 {
		return wtxmgr.BlockMeta{Block: wtxmgr.Block{Height: -1}}
	}

	b := wtxmgr.BlockMeta{
		Block: wtxmgr.Block{Height: height},
		Time:  timeNow(),
	}
	// Give it a fake block hash created from the height and time.
	binary.LittleEndian.PutUint32(b.Hash[0:4], uint32(height))
	binary.LittleEndian.PutUint64(b.Hash[4:12], uint64(b.Time.Unix()))
	return b
}

// Returns time.Now() with seconds resolution, this is what Store saves.
func timeNow() time.Time {
	return time.Unix(time.Now().Unix(), 0)
}
