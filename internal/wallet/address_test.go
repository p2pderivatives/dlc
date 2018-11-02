package wallet

import (
	"fmt"
	"testing"

	"github.com/btcsuite/btcwallet/waddrmgr"
	"github.com/stretchr/testify/assert"
)

func TestNewPubkey(t *testing.T) {
	tearDownFunc, wallet := setupWallet(t)
	defer tearDownFunc()

	assert.NotNil(t, wallet)
	fmt.Printf("WALLET\n%+v\n", wallet)

	wallet.CreateAccount(waddrmgr.KeyScopeBIP0084, testAccountName, testPrivPass)

	pub, _ := wallet.NewPubkey()

	assert.NotNil(t, pub)
}

func TestWitnessNewPubkeyScript(t *testing.T) {
}
