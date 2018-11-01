package wallet

import (
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcwallet/waddrmgr"
	"github.com/btcsuite/btcwallet/walletdb"
	_ "github.com/btcsuite/btcwallet/walletdb/bdb" // blank import for bolt db driver
	"github.com/btcsuite/btcwallet/wtxmgr"
)

// Namespace bucket keys.
var (
	waddrmgrNamespaceKey = []byte("waddrmgr")
	wtxmgrNamespaceKey   = []byte("wtxmgr")
)

// Wallet is hierarchical deterministic wallet
type Wallet struct {
	params           *chaincfg.Params
	publicPassphrase []byte // I'm thinking this should removed...
	// rpc    *rpc.BtcRPC

	db      walletdb.DB
	Manager *waddrmgr.Manager
	TxStore *wtxmgr.Store
}

// CreateWallet returns a new Wallet, also creates db where wallet resides
// TODO: separate db creation and Manager creation
// TODO: create loader script for wallet init
func CreateWallet(params *chaincfg.Params, seed, pubPass, privPass []byte,
	dbFilePath, walletName string) (*Wallet, error) {
	// TODO: add prompts for dbDirPath, walletDBname
	// Create a new db at specified path
	dbDirPath := filepath.Join(dbFilePath, params.Name)
	dbPath := filepath.Join(dbDirPath, walletName+".db")
	exists, err := fileExists(dbPath)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("something already exists on this filepath")
	}
	err = os.MkdirAll(dbDirPath, 0700)
	if err != nil {
		return nil, err
	}

	db, err := walletdb.Create("bdb", dbPath)
	if err != nil {
		return nil, err
	}

	// Create Wallet struct
	err = Create(db, params, seed, pubPass, privPass)
	if err != nil {
		return nil, err
	}

	// Open the wallet
	wallet, err := Open(db, pubPass, params)

	return wallet, err
}

// Create creates an new wallet, writing it to the passed in db.
func Create(db walletdb.DB, params *chaincfg.Params, seed, pubPass,
	privPass []byte) error {
	err := walletdb.Update(db, func(tx walletdb.ReadWriteTx) error {
		addrmgrNs, e := tx.CreateTopLevelBucket(waddrmgrNamespaceKey)
		if e != nil {
			return e
		}
		txmgrNs, e := tx.CreateTopLevelBucket(wtxmgrNamespaceKey)
		if e != nil {
			return e
		}

		birthday := time.Now()
		e = waddrmgr.Create(
			addrmgrNs, seed, pubPass, privPass, params, nil,
			birthday,
		)
		if e != nil {
			// TODO: figure out how to gracefully close db
			//   possibly defer db.Close() ?
			db.Close()
			return e
		}
		e = wtxmgr.Create(txmgrNs)
		return e
	})

	return err
}

// Open loads a wallet from the passed db and public pass phrase.
func Open(db walletdb.DB, pubPass []byte,
	params *chaincfg.Params) (*Wallet, error) {
	err := walletdb.View(db, func(tx walletdb.ReadTx) error {
		waddrmgrBucket := tx.ReadBucket(waddrmgrNamespaceKey)
		if waddrmgrBucket == nil {
			return errors.New("missing address manager namespace")
		}
		wtxmgrBucket := tx.ReadBucket(wtxmgrNamespaceKey)
		if wtxmgrBucket == nil {
			return errors.New("missing transaction manager namespace")
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// TODO: Perform wallet upgrades if necessary?

	// Open database abstraction instances
	var (
		addrMgr *waddrmgr.Manager
		txMgr   *wtxmgr.Store
	)
	err = walletdb.View(db, func(tx walletdb.ReadTx) error {
		addrmgrNs := tx.ReadBucket(waddrmgrNamespaceKey)
		txmgrNs := tx.ReadBucket(wtxmgrNamespaceKey)
		var e error
		addrMgr, e = waddrmgr.Open(addrmgrNs, pubPass, params)
		if e != nil {
			return e
		}
		txMgr, e = wtxmgr.Open(txmgrNs, params)
		return e
	})
	if err != nil {
		return nil, err
	}

	w := &Wallet{
		params:           params,
		publicPassphrase: pubPass,
		db:               db,
		Manager:          addrMgr,
		TxStore:          txMgr,
	}

	return w, nil
}

// TODO: add Close wallet function that will gracefully close db, Manager...

// CreateAccount creates a new account in ScopedKeyManagar of scope
func (w *Wallet) CreateAccount(scope waddrmgr.KeyScope, name string,
	privPass []byte) (uint32, error) {
	// unlock Manager
	err := walletdb.Update(w.db, func(tx walletdb.ReadWriteTx) error {
		ns := tx.ReadWriteBucket(waddrmgrNamespaceKey)
		e := w.Manager.Unlock(ns, privPass)
		return e
	})
	if err != nil {
		return 0, err
	}

	scopedMgr, err := w.Manager.FetchScopedKeyManager(scope)
	if err != nil {
		return 0, err
	}

	var account uint32
	err = walletdb.Update(w.db, func(tx walletdb.ReadWriteTx) error {
		ns := tx.ReadWriteBucket(waddrmgrNamespaceKey)
		var e error
		account, e = scopedMgr.NewAccount(ns, name)
		return e
	})
	if err != nil {
		return 0, err
	}

	return account, nil
}

// NewAddress returns a new ManagedAddress for a given scope and account number.
// NOTE: this function callsNextExternalAddresses to generate a ManagadAdddress.
func (w *Wallet) NewAddress(scope waddrmgr.KeyScope, privPass []byte,
	account uint32, numAddresses uint32) ([]waddrmgr.ManagedAddress, error) {
	// unlock Manager
	err := walletdb.Update(w.db, func(tx walletdb.ReadWriteTx) error {
		ns := tx.ReadWriteBucket(waddrmgrNamespaceKey)
		e := w.Manager.Unlock(ns, privPass)
		return e
	})
	if err != nil {
		return nil, err
	}

	// get ScopedKeyManager
	scopedMgr, err := w.Manager.FetchScopedKeyManager(scope)
	if err != nil {
		return nil, err
	}

	var addrs []waddrmgr.ManagedAddress
	err = walletdb.Update(w.db, func(tx walletdb.ReadWriteTx) error {
		ns := tx.ReadWriteBucket(waddrmgrNamespaceKey)
		var e error
		addrs, e = scopedMgr.NextExternalAddresses(ns, account, numAddresses)
		return e
	})
	if err != nil {
		return nil, err
	}

	return addrs, nil
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
