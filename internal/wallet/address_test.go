package wallet

import (
	"testing"

	"github.com/btcsuite/btcwallet/waddrmgr"
	"github.com/stretchr/testify/assert"
)

func TestNewPubkey(t *testing.T) {
	tearDownFunc, wallet := setupWallet(t)
	defer tearDownFunc()

	wallet.CreateAccount(waddrmgr.KeyScopeBIP0084, testAccountName, testPrivPass)
	pub, _ := wallet.NewPubkey()

	assert.NotNil(t, pub)
}

func TestWitnessNewPubkeyScript(t *testing.T) {
}
