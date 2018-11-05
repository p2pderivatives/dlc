package wallet

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPubkey(t *testing.T) {
	wallet, tearDownFunc := setupWallet(t)
	defer tearDownFunc()

	pub, err := wallet.NewPubkey()

	assert.Nil(t, err)
	assert.NotNil(t, pub)
}

func TestNewWitnessPubkeyScript(t *testing.T) {
	wallet, tearDownFunc := setupWallet(t)
	defer tearDownFunc()

	pkScript, err := wallet.NewWitnessPubkeyScript()

	assert.Nil(t, err)
	assert.NotEmpty(t, pkScript)
}
