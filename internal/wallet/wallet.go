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
	_ "github.com/btcsuite/btcwallet/walletdb/bdb" // blank import for bolt db driver
)

// Namespace bucket keys.
var (
	waddrmgrNamespaceKey = []byte("waddrmgr")
	wtxmgrNamespaceKey   = []byte("wtxmgr")
)

// Wallet is hierarchical deterministic wallet
type Wallet struct {
	params chaincfg.Params
	// rpc    *rpc.BtcRPC

	db               walletdb.DB
	Manager          *waddrmgr.Manager
	publicPassphrase []byte
}

// PublicKeyInfo is publickey data.
type PublicKeyInfo struct {
	idx uint32
	pub *btcec.PublicKey
	adr string
}

// NewWallet returns a new Wallet
func CreateWallet(params chaincfg.Params, seed, pubPass, privPass []byte, dbFilePath, walletName string) (*Wallet, error) {
	wallet := &Wallet{}
	wallet.params = params
	// wallet.rpc = rpc
	wallet.publicPassphrase = pubPass

	// TODO: add prompts for dbDirPath, walletDBname
	dbDirPath := filepath.Join(dbFilePath, params.Name)
	walletDBname := walletName + ".db"
	dbPath := filepath.Join(dbDirPath, walletDBname)
	exists, err := fileExists(dbPath)
	if err != nil {
		return nil, err
	}
	if exists {
		fmt.Printf("Something already exists on this filepath!")
		return nil, err
	}
	err = os.MkdirAll(dbDirPath, 0700)
	if err != nil {
		return nil, err
	}

	db, err := walletdb.Create("bdb", dbPath)
	if err != nil {
		_ = os.RemoveAll(dbDirPath)
		fmt.Println(err)
		return nil, err
	}
	wallet.db = db

	var mgr *waddrmgr.Manager
	err = walletdb.Update(db, func(tx walletdb.ReadWriteTx) error {
		addrmgrNs, err := tx.CreateTopLevelBucket(waddrmgrNamespaceKey)
		if err != nil {
			return err
		}
		// TODO: figure out if txmgrNs is needed
		//txmgrNs, err := tx.CreateTopLevelBucket(wtxmgrNamespaceKey)
		//if err != nil {
		//	return err
		//}

		birthday := time.Now()
		err = waddrmgr.Create(
			addrmgrNs, seed, pubPass, privPass, &params, nil,
			birthday,
		)
		if err != nil {
			db.Close()
			return err
		}
		mgr, err = waddrmgr.Open(addrmgrNs, pubPass, &params)
		wallet.Manager = mgr

		return err
	})

	return wallet, nil
}

func (w *Wallet) CreateAccount(scope waddrmgr.KeyScope, name string, privPass []byte) (uint32, error) {
	// unlock Manager
	err := walletdb.Update(w.db, func(tx walletdb.ReadWriteTx) error {
		ns := tx.ReadWriteBucket(waddrmgrNamespaceKey)
		err := w.Manager.Unlock(ns, privPass)
		return err
	})
	if err != nil {
		return 0, err
	}

	scopedMgr, err := w.Manager.FetchScopedKeyManager(scope)

	var account uint32
	err = walletdb.Update(w.db, func(tx walletdb.ReadWriteTx) error {
		ns := tx.ReadWriteBucket(waddrmgrNamespaceKey)
		var err error
		account, err = scopedMgr.NewAccount(ns, name)
		return err
	})
	if err != nil {
		fmt.Printf("NewAccount: unexpected error: %v", err)
		return 0, err
	}

	return account, err
}

// Helper function, TODO: move somewhere else?
func fileExists(filePath string) (bool, error) {
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
