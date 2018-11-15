package wallet

import (
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcwallet/waddrmgr"
	"github.com/btcsuite/btcwallet/walletdb"
	_ "github.com/btcsuite/btcwallet/walletdb/bdb" // blank import for bolt db driver
	"github.com/dgarage/dlc/internal/rpc"
)

// Wallet is an interface that provides access to manage pubkey addresses and
// sign scripts of managed addressesc using private key. It also manags utxos.
type Wallet interface {
	NewPubkey() (*btcec.PublicKey, error)

	// WitnessSignature returns witness signature for a given txin and pubkey
	WitnessSignature(
		tx *wire.MsgTx, idx int, amt btcutil.Amount, sc []byte, pub *btcec.PublicKey,
	) (sign []byte, err error)

	// WitnessSignatureWithCallback does the same with WitnessSignature do
	// applying a given func to private key before calculating signature
	WitnessSignatureWithCallback(
		tx *wire.MsgTx, idx int, amt btcutil.Amount, sc []byte, pub *btcec.PublicKey,
		privkeyConverter PrivateKeyConverter,
	) (sign []byte, err error)

	ListUnspent() (utxos []Utxo, err error)

	// SelectUtxos selects utxos for requested amount
	// by considering additional fee per txin and txout
	SelectUnspent(
		amt, feePerTxIn, feePerTxOut btcutil.Amount,
	) (utxos []Utxo, change btcutil.Amount, err error)

	// Unlock unlocks address manager
	Unlock(privPass []byte) error

	Close() error
}

// PrivateKeyConverter is a callback func applied to private key before creating witness signature
type PrivateKeyConverter func(*btcec.PrivateKey) (*btcec.PrivateKey, error)

// Namespace bucket keys.
var (
	waddrmgrNamespaceKey = []byte("waddrmgr")
	waddrmgrKeyScope     = waddrmgr.KeyScopeBIP0084

	// TODO: have rpc params be read from conf file?
	rpcport     = "localhost: REPLACEME"
	rpcusername = "RENAME!"
	rpcpassword = "RENAME!"
)

const accountName = "dlc"

// Wallet is hierarchical deterministic wallet
type wallet struct {
	params           *chaincfg.Params
	publicPassphrase []byte
	rpc              rpc.Client
	db               walletdb.DB
	manager          *waddrmgr.Manager
	account          uint32
}

// wallet should satisfy Wallet interface
var _ Wallet = (*wallet)(nil)

// CreateWallet returns a new Wallet, also creates db where wallet resides
// TODO: separate db creation and Manager creation
// TODO: create loader script for wallet init
// TODO: add prompts for dbDirPath, walletDBname
func CreateWallet(
	params *chaincfg.Params,
	seed, pubPass, privPass []byte,
	dbFilePath, walletName string) (Wallet, error) {

	dbDirPath := filepath.Join(dbFilePath, params.Name)
	db, err := createDB(dbDirPath, walletName+".db")
	if err != nil {
		return nil, err
	}

	return create(db, params, seed, pubPass, privPass)
}

// createDB creates a new db at specified path
func createDB(dbDirPath, dbname string) (walletdb.DB, error) {
	dbPath := filepath.Join(dbDirPath, dbname)
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

	return db, err
}

// Create creates an new wallet, writing it to the passed in db.
func create(
	db walletdb.DB,
	params *chaincfg.Params,
	seed, pubPass, privPass []byte) (*wallet, error) {

	err := createManagers(db, seed, pubPass, privPass, params)
	if err != nil {
		return nil, err
	}

	err = createAccount(db, privPass, pubPass, params)
	if err != nil {
		return nil, err
	}

	w, err := open(db, pubPass, params)
	if err != nil {
		return nil, err
	}

	return w, nil
}

// createManagers create address manager and tx manager
func createManagers(
	db walletdb.DB,
	seed, pubPass, privPass []byte,
	params *chaincfg.Params,
) error {
	return walletdb.Update(db, func(tx walletdb.ReadWriteTx) error {
		addrmgrNs, e := tx.CreateTopLevelBucket(waddrmgrNamespaceKey)
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
		return e
	})
}

// createAccount creates a new account in ScopedKeyManagar of scope
func createAccount(
	db walletdb.DB, privPass, pubPass []byte, params *chaincfg.Params) error {
	return walletdb.Update(db, func(tx walletdb.ReadWriteTx) error {
		ns := tx.ReadWriteBucket(waddrmgrNamespaceKey)
		addrMgr, e := waddrmgr.Open(ns, pubPass, params)
		if e != nil {
			return e
		}

		e = addrMgr.Unlock(ns, privPass)
		if e != nil {
			return e
		}

		scopedMgr, e := addrMgr.FetchScopedKeyManager(waddrmgrKeyScope)
		if e != nil {
			return e
		}

		_, e = scopedMgr.NewAccount(ns, accountName)
		return e
	})
}

// Open loads a wallet from the passed db and public pass phrase.
func Open(
	db walletdb.DB, pubPass []byte, params *chaincfg.Params,
) (Wallet, error) {
	return open(db, pubPass, params)
}

// open is an implementation of Open
func open(
	db walletdb.DB, pubPass []byte, params *chaincfg.Params,
) (*wallet, error) {
	// TODO: Perform wallet upgrades/updates if necessary?

	// Open database abstraction instances
	var (
		addrMgr *waddrmgr.Manager
		account uint32
	)
	err := walletdb.View(db, func(tx walletdb.ReadTx) error {
		addrmgrNs := tx.ReadBucket(waddrmgrNamespaceKey)
		if addrmgrNs == nil {
			return errors.New("missing address manager namespace")
		}

		var e error
		addrMgr, e = waddrmgr.Open(addrmgrNs, pubPass, params)
		if e != nil {
			return e
		}
		scopedMgr, e := addrMgr.FetchScopedKeyManager(waddrmgrKeyScope)
		if e != nil {
			return e
		}
		account, e = scopedMgr.LookupAccount(addrmgrNs, accountName)
		if e != nil {
			return e
		}

		return e
	})
	if err != nil {
		return nil, err
	}

	rpc, err := rpc.NewClient(rpcport, rpcusername, rpcpassword)
	if err != nil {
		return nil, err
	}

	w := &wallet{
		params:           params,
		publicPassphrase: pubPass,
		rpc:              rpc,
		db:               db,
		manager:          addrMgr,
		account:          account,
	}

	return w, nil
}

// Unlock unlocks address manager with a given private pass phrase
func (w *wallet) Unlock(privPass []byte) error {
	return walletdb.Update(w.db, func(tx walletdb.ReadWriteTx) error {
		ns := tx.ReadWriteBucket(waddrmgrNamespaceKey)
		return w.manager.Unlock(ns, privPass)
	})
}

// Close closes managers
func (w *wallet) Close() error {
	w.manager.Close()
	return nil
}

// Helper function
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
