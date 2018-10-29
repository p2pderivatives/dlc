package wallet

import (
	"fmt"
	"os"
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/btcsuite/btcwallet/waddrmgr"
	"github.com/stretchr/testify/assert"
)

func TestNewWallet(t *testing.T) {
	params := chaincfg.RegressionNetParams
	pubPass := []byte("testpub")
	privPass := []byte("testpri")
	dbFilePath := "~/testdb"
	seed, _ := hdkeychain.GenerateSeed(hdkeychain.RecommendedSeedLen)

	wallet, _ := NewWallet(params, seed, pubPass, privPass, dbFilePath)
	assert.NotNil(t, wallet)
	assert.NotNil(t, wallet.Manager)

	// delete created db
	_ = os.RemoveAll(dbFilePath)
}

// TODO: create testing interface

func TestCreateNewAccount(t *testing.T) {
	params := chaincfg.RegressionNetParams
	pubPass := []byte("testpub")
	privPass := []byte("testpri")
	dbFilePath := "~/testdb2"
	seed, _ := hdkeychain.GenerateSeed(hdkeychain.RecommendedSeedLen)

	wallet, _ := NewWallet(params, seed, pubPass, privPass, dbFilePath)
	assert.NotNil(t, wallet.Manager)

	testAccountName := "testy"
	account, _ := wallet.NewAccount(waddrmgr.KeyScopeBIP0084, testAccountName, privPass)

	// TODO: actually use a testing library
	fmt.Printf("account: %d, \n", account)

	_ = os.RemoveAll(dbFilePath)
}
