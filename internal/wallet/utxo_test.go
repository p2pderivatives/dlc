package wallet

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListUnspent(t *testing.T) {
	// dbFilePath := "./testdb"

	// Create a temporary directory for testing.
	dirName, err := ioutil.TempDir("", "managertest")
	if err != nil {
		t.Fatalf("Failed to create db temp dir: %v", err)
	}

	wallet, err := CreateWallet(&params, seed, pubPassphrase, privPassphrase,
		dirName, walletName)

	assert.Nil(t, err)
	assert.NotNil(t, wallet.publicPassphrase)
	assert.NotNil(t, wallet.params)
	assert.NotNil(t, wallet.Manager)
	assert.NotNil(t, wallet.TxStore)

	// delete created db
	_ = os.RemoveAll(dirName)
}
