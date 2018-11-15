package wallet

import (
	"testing"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/dgarage/dlc/internal/mocks/rpcmock"
	"github.com/stretchr/testify/assert"
)

// Test setup?
// 		create wallet
// 		mine regtest coins?
//  	ListUnspent() to check if we can see the mined coins?

// TestListUnspent() will also need to check different types of scripts
func TestListUnspent(t *testing.T) {
	w, tearDownFunc := setupWallet(t)
	defer tearDownFunc()

	var utxos []btcjson.ListUnspentResult
	utxo := btcjson.ListUnspentResult{
		TxID:          "ce9d930c2664547ad8aba6944c8047321bde0c1c1d6551c41ebb8d9ad975dd0b",
		Vout:          uint32(0),
		Address:       "tb1qds49lkplvws9q4df04e5e9nq5d6asnkkhna8hg",
		Account:       "",
		ScriptPubKey:  "00146c2a5fd83f63a05055a97d734c9660a375d84ed6",
		RedeemScript:  "",
		Amount:        float64(0.31864472),
		Confirmations: int64(30006),
		Spendable:     true,
	}
	utxos = append(utxos, utxo)

	rpcc := &rpcmock.Client{}
	rpcc = mockListUnspent(rpcc, utxos, nil)
	w.rpc = rpcc

	rutxos, err := w.ListUnspent()

	assert.Nil(t, err)
	assert.Equal(t, rutxos, utxos)
}

func mockListUnspent(c *rpcmock.Client, utxos []btcjson.ListUnspentResult, err error) *rpcmock.Client {
	c.On("ListUnspent").Return(utxos, err)

	return c
}
