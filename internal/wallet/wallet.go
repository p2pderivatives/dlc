// Package wallet project wallet.go
package wallet

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcwallet/waddrmgr"
	"github.com/btcsuite/btcwallet/walletdb"
	_ "github.com/btcsuite/btcwallet/walletdb/bdb"
	"github.com/btcsuite/btcwallet/wtxmgr"
)

const (
	PUBPASSPHRASE = "pubpass"
	PRIPASSPHRASE = "pripass"
	// BIRTHDAY  Time.time =
)

// Namespace bucket keys.
var (
	waddrmgrNamespaceKey = []byte("waddrmgr")
	wtxmgrNamespaceKey   = []byte("wtxmgr")
)

// Wallet is hierarchical deterministic wallet
type Wallet struct {
	params chaincfg.Params
	size   int // is this parameter still needed?
	// rpc    *rpc.BtcRPC

	db      walletdb.DB
	Manager *waddrmgr.Manager
}

// PublicKeyInfo is publickey data.
type PublicKeyInfo struct {
	idx uint32
	pub *btcec.PublicKey
	adr string
}

// NewWallet returns a new Wallet
// func NewWallet(params chaincfg.Params, rpc *rpc.BtcRPC, seed []byte) (*Wallet, error) {
func NewWallet(params chaincfg.Params, seed []byte) (*Wallet, error) {
	wallet := &Wallet{}
	wallet.params = params
	// wallet.rpc = rpc
	wallet.size = 16

	dbPath := filepath.Join(os.TempDir(), "dev.db")
	db, err := walletdb.Create("bdb", dbPath)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	err = walletdb.Update(db, func(tx walletdb.ReadWriteTx) error {
		addrmgrNs, err := tx.CreateTopLevelBucket(waddrmgrNamespaceKey)
		if err != nil {
			return err
		}
		txmgrNs, err := tx.CreateTopLevelBucket(wtxmgrNamespaceKey)
		if err != nil {
			return err
		}

		birthday := time.Now()
		err = waddrmgr.Create(
			addrmgrNs, seed, []byte(PUBPASSPHRASE), []byte(PRIPASSPHRASE), &params, nil,
			birthday,
		)
		if err != nil {
			return err
		}
		return wtxmgr.Create(txmgrNs)
	})
	return wallet, nil
}

// func (w *Wallet) SendTx(tx *wire.MsgTx) (*chainhash.Hash, error) {
// 	allowHighFees := false
// 	//return w.rpc.SendRawTransaction(tx, allowHighFees)

// 	// testing
// 	// marshalled: `{"jsonrpc":"1.0","method":"sendrawtransaction","params":["1122"],"id":1}`,
// 	// unmarshalled: &btcjson.SendRawTransactionCmd{
// 	// 	HexTx:         "1122",
// 	// 	AllowHighFees: btcjson.Bool(false),
// 	// },
// 	// https://github.com/btcsuite/btcd/blob/fdfc19097e7ac6b57035062056f5b7b4638b8898/btcjson/chainsvrcmds_test.go#L903

// }
