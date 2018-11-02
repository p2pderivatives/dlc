package wallet

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/btcsuite/btcwallet/waddrmgr"
	"github.com/stretchr/testify/assert"
)

func TestCreateWallet(t *testing.T) {
	dirName, _ := ioutil.TempDir("", "testcreatewallet")
	defer os.RemoveAll(dirName)

	w, err := CreateWallet(
		testNetParams, testSeed, testPubPass, testPrivPass, dirName, testWalletName)

	assert := assert.New(t)
	assert.Nil(err)
	assert.NotNil(w)
}

// TODO: add tests for Create(...) and Open(...)

func TestCreateAccount(t *testing.T) {
	tearDownFunc, wallet := setupWallet(t)
	defer tearDownFunc()

	expectedAccountNumber := uint32(1)

	account, _ := wallet.CreateAccount(
		waddrmgr.KeyScopeBIP0084, testAccountName, testPrivPass)

	assert.Equal(t, expectedAccountNumber, account)
}
