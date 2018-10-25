package wallet

import (
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/stretchr/testify/assert"
)

func TestNewWallet(t *testing.T) {
	params := chaincfg.RegressionNetParams
	pubPass := []byte("testpub")
	privPass := []byte("testpri")
	seed, _ := hdkeychain.GenerateSeed(hdkeychain.RecommendedSeedLen)

	wallet, _ := NewWallet(params, seed, pubPass, privPass)
	assert.NotNil(t, wallet)
}
