package wallet

import (
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcwallet/waddrmgr"
	"github.com/btcsuite/btcwallet/walletdb"
	_ "github.com/btcsuite/btcwallet/walletdb/bdb" // blank import for bolt db driver
	"github.com/btcsuite/btcwallet/wtxmgr"
)

// Wallet is an interface that provides access to manage pubkey addresses and
// sign scripts of managed addressesc using private key. It also manags utxos.
type Wallet interface {
	CreateAccount(
		scope waddrmgr.KeyScope, name string, privPass []byte,
	) (account uint32, err error)

	NewPubkey() (*btcec.PublicKey, error)
	NewWitnessPubkeyScript() (pkScript []byte, err error)
	ListUnspent() (utxos []Utxo, err error)

	Close() error
}

// Namespace bucket keys.
var (
	waddrmgrNamespaceKey = []byte("waddrmgr")
	wtxmgrNamespaceKey   = []byte("wtxmgr")
)

// Wallet is hierarchical deterministic wallet
type wallet struct {
	params            *chaincfg.Params
	publicPassphrase  []byte
	privatePassphrase []byte
	// rpc    *rpc.BtcRPC

	db      walletdb.DB
	manager *waddrmgr.Manager
	txStore *wtxmgr.Store
}

// CreateWallet returns a new Wallet, also creates db where wallet resides
// TODO: separate db creation and Manager creation
// TODO: create loader script for wallet init
func CreateWallet(params *chaincfg.Params, seed, pubPass, privPass []byte,
	dbFilePath, walletName string) (*wallet, error) {
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
	return Open(db, pubPass, privPass, params)
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
			return e
		}
		e = wtxmgr.Create(txmgrNs)
		return e
	})

	return err
}

// Open loads a wallet from the passed db and public pass phrase.
func Open(db walletdb.DB, pubPass, privPass []byte,
	params *chaincfg.Params) (*wallet, error) {
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

	// TODO: Perform wallet upgrades/updates if necessary?

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

	w := &wallet{
		params:            params,
		publicPassphrase:  pubPass,
		privatePassphrase: privPass,
		db:                db,
		manager:           addrMgr,
		txStore:           txMgr,
	}

	return w, nil
}

// Close closes managers
func (w *wallet) Close() error {
	w.manager.Close()
	return nil
}

// CreateAccount creates a new account in ScopedKeyManagar of scope
func (w *wallet) CreateAccount(scope waddrmgr.KeyScope, name string,
	privPass []byte) (uint32, error) {
	// unlock Manager
	err := walletdb.Update(w.db, func(tx walletdb.ReadWriteTx) error {
		ns := tx.ReadWriteBucket(waddrmgrNamespaceKey)
		e := w.manager.Unlock(ns, privPass)
		return e
	})
	if err != nil {
		return 0, err
	}

	scopedMgr, err := w.manager.FetchScopedKeyManager(scope)
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
