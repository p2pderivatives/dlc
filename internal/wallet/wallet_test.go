package wallet

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/btcsuite/btcwallet/waddrmgr"
	"github.com/stretchr/testify/assert"
)

var (
	pubPassphrase  = []byte("_DJr{fL4H0O}*-0\n:V1izc)(6BomK")
	privPassphrase = []byte("81lUHXnOMZ@?XXd7O9xyDIWIbXX-lj")

	//dirName    = "./testdb"
	walletName = "testwallet"

	params = chaincfg.RegressionNetParams

	waddrmgrTestNamespaceKey = []byte("waddrmgrNamespace")
)

func TestCreateWallet(t *testing.T) {
	params := chaincfg.RegressionNetParams
	pubPass := []byte("testpub")
	privPass := []byte("testpri")
	dbFilePath := "./testdb"
	walletName := "testwallet"
	seed, _ := hdkeychain.GenerateSeed(hdkeychain.RecommendedSeedLen)

	wallet, _ := CreateWallet(params, seed, pubPass, privPass, dbFilePath, walletName)
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
func setupManager(t *testing.T) (tearDownFunc func(), wallet *Wallet) {
	// Create a temporary directory for testing.
	dirName, err := ioutil.TempDir("", "managertest")
	if err != nil {
		t.Fatalf("Failed to create db temp dir: %v", err)
	}

	seed, _ := hdkeychain.GenerateSeed(hdkeychain.RecommendedSeedLen)

	wallet, err = CreateWallet(params, seed, pubPassphrase, privPassphrase, dirName, walletName)
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
	tearDownFunc, wallet := setupManager(t)
	defer tearDownFunc()

	expectedAccountNumber := uint32(1)

	testAccountName := "testy"
	account, _ := wallet.CreateAccount(waddrmgr.KeyScopeBIP0084, testAccountName, privPassphrase)

	assert.Equal(t, expectedAccountNumber, account)
}

func TestNewExternalAddress(t *testing.T) {
	tearDownFunc, wallet := setupManager(t)
	defer tearDownFunc()

	testAccountName := "testy"
	account, _ := wallet.CreateAccount(waddrmgr.KeyScopeBIP0084, testAccountName, privPassphrase)

	numAddresses := uint32(4)

	addrs, _ := wallet.NewExternalAddress(waddrmgr.KeyScopeBIP0084, privPassphrase, account, numAddresses)
	fmt.Printf("%+v", addrs)
	fmt.Printf("%+v", addrs[0])
}
