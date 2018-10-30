package wallet

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcwallet/waddrmgr"
	"github.com/btcsuite/btcwallet/walletdb"
	"github.com/stretchr/testify/assert"
)

var (
	seed = []byte{
		0xa7, 0x97, 0x63, 0xcf, 0x88, 0x54, 0xe1, 0xd3, 0xb0,
		0x89, 0x07, 0xed, 0xc6, 0x96, 0x05, 0xf3, 0x38, 0xc1,
		0xb6, 0xb8, 0x39, 0xbe, 0xd9, 0xfd, 0x21, 0x6a, 0x6c,
		0x03, 0xce, 0xe2, 0x2c, 0x84,
	}

	params         = chaincfg.RegressionNetParams
	pubPassphrase  = []byte("_DJr{fL4H0O}*-0\n:V1izc)(6BomK")
	privPassphrase = []byte("81lUHXnOMZ@?XXd7O9xyDIWIbXX-lj")

	//dirName    = "./testdb"
	walletName               = "testwallet"
	waddrmgrTestNamespaceKey = []byte("waddrmgr")
)

func TestCreateWallet(t *testing.T) {
	params := chaincfg.RegressionNetParams
	dbFilePath := "./testdb"

	wallet, _ := CreateWallet(params, seed, pubPassphrase, privPassphrase,
		dbFilePath, walletName)
	assert.NotNil(t, wallet)
	assert.NotNil(t, wallet.Manager)
	assert.NotNil(t, wallet.db)
	assert.NotNil(t, wallet.publicPassphrase)

	// delete created db
	_ = os.RemoveAll(dbFilePath)
}

// TODO: create testing interface
// setupManager creates a new address manager and returns a teardown function
// that should be invoked to ensure it is closed and removed upon completion.
func setupWallet(t *testing.T) (tearDownFunc func(), wallet *Wallet) {
	// Create a temporary directory for testing.
	dirName, err := ioutil.TempDir("", "managertest")
	if err != nil {
		t.Fatalf("Failed to create db temp dir: %v", err)
	}

	wallet, err = CreateWallet(params, seed, pubPassphrase, privPassphrase,
		dirName, walletName)
	if err != nil {
		wallet.db.Close()
		_ = os.RemoveAll(dirName)
		t.Fatalf("Failed to create test Wallet: %v", err)
	}

	tearDownFunc = func() {
		wallet.Manager.Close()
		wallet.db.Close()
		_ = os.RemoveAll(dirName)
	}

	return tearDownFunc, wallet
}

func TestCreateAccount(t *testing.T) {
	tearDownFunc, wallet := setupWallet(t)
	defer tearDownFunc()

	expectedAccountNumber := uint32(1)

	testAccountName := "testy"
	account, _ := wallet.CreateAccount(waddrmgr.KeyScopeBIP0084, testAccountName,
		privPassphrase)

	assert.Equal(t, expectedAccountNumber, account)
}

func TestNewAddress(t *testing.T) {
	tearDownFunc, wallet := setupWallet(t)
	defer tearDownFunc()

	testAccountName := "testy"
	account, _ := wallet.CreateAccount(waddrmgr.KeyScopeBIP0084, testAccountName,
		privPassphrase)

	numAddresses := uint32(1)

	addrs, _ := wallet.NewAddress(waddrmgr.KeyScopeBIP0084,
		privPassphrase, account, numAddresses)
	err := walletdb.View(wallet.db, func(tx walletdb.ReadTx) error {
		ns := tx.ReadBucket(waddrmgrTestNamespaceKey)
		assert.False(t, addrs[0].Used(ns))
		return nil
	})
	if err != nil {
		t.Errorf("Unlock: unexpected error: %v", err)
	}
}
