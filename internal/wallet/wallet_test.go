package wallet

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/stretchr/testify/assert"
)

func TestCreateWallet(t *testing.T) {
	assert := assert.New(t)
	dirName, _ := ioutil.TempDir("", "testcreatewallet")
	defer os.RemoveAll(dirName)

	seed, err := hdkeychain.GenerateSeed(
		hdkeychain.RecommendedSeedLen)
	assert.NoError(err)

	w, err := CreateWallet(
		testNetParams, seed, testPubPass, testPrivPass, dirName, testWalletName)

	// assertions
	assert.Nil(err)
	_, ok := w.(Wallet)
	assert.True(ok)
}

func TestOpen(t *testing.T) {
	_w, tearDownFunc := setupWallet(t)
	defer tearDownFunc()

	// close wallet
	_w.Close()

	// open wallet
	db := _w.db
	pubPass := _w.publicPassphrase
	params := _w.params
	w, err := open(db, pubPass, params, nil)

	// assertions
	assert := assert.New(t)
	assert.Nil(err)

	// test if the oepned account is the same with the created one
	assert.Equal(_w.account, w.account)

	// test if it satisfies Wallet interface
	var W Wallet = w
	_, ok := W.(Wallet)
	assert.True(ok)
}
