package wallet

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/btcsuite/btcwallet/walletdb"
	"github.com/stretchr/testify/assert"
)

var (
	testNetParams  = &chaincfg.RegressionNetParams
	testPubPass    = []byte("_DJr{fL4H0O}*-0\n:V1izc)(6BomK")
	testPrivPass   = []byte("81lUHXnOMZ@?XXd7O9xyDIWIbXX-lj")
	testWalletName = "testwallet.db"
)

func setupDB(t *testing.T) (db walletdb.DB, tearDownFunc func()) {
	assert := assert.New(t)

	dbDirPath, err := ioutil.TempDir("", "testdb")
	assert.Nil(err)

	db, err = createDB(dbDirPath, testWalletName)
	assert.Nil(err)

	tearDownFunc = func() {
		err = db.Close()
		assert.Nil(err)

		err = os.RemoveAll(dbDirPath)
		assert.Nil(err)
	}
	return
}

func setupWallet(t *testing.T) (*wallet, func()) {
	assert := assert.New(t)
	db, deleteDB := setupDB(t)

	seed, err := hdkeychain.GenerateSeed(
		hdkeychain.RecommendedSeedLen)
	assert.NoError(err)
	w, err := create(db, testNetParams, seed, testPubPass, testPrivPass)
	assert.NoError(err)

	tearDownFunc := func() {
		err = w.Close()
		assert.Nil(err)
		deleteDB()
	}

	return w, tearDownFunc
}
