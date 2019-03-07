package wallet

import (
	"testing"

	"github.com/btcsuite/btcutil"
	"github.com/p2pderivatives/dlc/internal/rpc"
	"github.com/stretchr/testify/assert"
)

// TestListUnspent() will also need to check different types of scripts
func TestListUnspent(t *testing.T) {
	assert := assert.New(t)
	w, tearDownFunc := setupWallet(t)
	defer tearDownFunc()

	c, err := rpc.NewTestRPCClient()
	assert.NoError(err)
	w.SetRPCClient(c)

	// empty at beginning
	utxos, err := w.ListUnspent()
	assert.NoError(err)
	assert.Empty(utxos)

	addr, err := w.NewAddress()
	assert.NoError(err)
	_ = rpc.Faucet(addr, 1*btcutil.SatoshiPerBitcoin)

	utxos, err = w.ListUnspent()
	assert.NoError(err)
	assert.Len(utxos, 1)
}
