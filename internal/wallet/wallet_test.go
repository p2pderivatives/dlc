package wallet

import (
	"os"
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/btcsuite/btcwallet/waddrmgr"
	"github.com/stretchr/testify/assert"
)

func TestCreateWallet(t *testing.T) {
	params := chaincfg.RegressionNetParams
	pubPass := []byte("testpub")
	privPass := []byte("testpri")
	dbFilePath := "./testdb"
	walletName := "testwallet"
	seed, _ := hdkeychain.GenerateSeed(hdkeychain.RecommendedSeedLen)

	wallet, _ := CreateWallet(params, seed, pubPass, privPass, dbFilePath, walletName)
	defer os.RemoveAll(dbFilePath)

	assert.NotNil(t, wallet)
	assert.NotNil(t, wallet.Manager)
	assert.NotNil(t, wallet.db)
	assert.NotNil(t, wallet.publicPassphrase)
}

// TODO: create testing interface

func TestCreateAccount(t *testing.T) {
	params := chaincfg.RegressionNetParams
	pubPass := []byte("testpub")
	privPass := []byte("testpri")
	dbFilePath := "./testdb2"
	walletName := "testwallet"

	seed, _ := hdkeychain.GenerateSeed(hdkeychain.RecommendedSeedLen)

	wallet, _ := CreateWallet(params, seed, pubPass, privPass, dbFilePath, walletName)
	defer os.RemoveAll(dbFilePath)

	assert.NotNil(t, wallet.Manager)

	expectedAccountNumber := uint32(1)

	testAccountName := "testy"
	account, err := wallet.CreateAccount(
		waddrmgr.KeyScopeBIP0084, testAccountName, privPass)

	assert := assert.New(t)
	assert.Nil(err)
	assert.Equal(expectedAccountNumber, account)
}
