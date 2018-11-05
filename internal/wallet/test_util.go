package wallet

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcwallet/walletdb"
	"github.com/stretchr/testify/assert"
)

var (
	testNetParams = &chaincfg.RegressionNetParams
	testSeed      = []byte{
		0xa7, 0x97, 0x63, 0xcf, 0x88, 0x54, 0xe1, 0xd3, 0xb0,
		0x89, 0x07, 0xed, 0xc6, 0x96, 0x05, 0xf3, 0x38, 0xc1,
		0xb6, 0xb8, 0x39, 0xbe, 0xd9, 0xfd, 0x21, 0x6a, 0x6c,
		0x03, 0xce, 0xe2, 0x2c, 0x84,
	}
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

	w, err := create(db, testNetParams, testSeed, testPubPass, testPrivPass)
	assert.Nil(err)

	tearDownFunc := func() {
		err = w.Close()
		assert.Nil(err)
		deleteDB()
	}

	return w, tearDownFunc
}
