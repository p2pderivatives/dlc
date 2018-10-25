package wallet

import (
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/stretchr/testify/assert"
)

func TestNewWallet(t *testing.T) {
	params := chaincfg.RegressionNetParams
	seed, _ := hdkeychain.GenerateSeed(hdkeychain.RecommendedSeedLen)

	wallet, _ := NewWallet(params, seed)
	assert.NotNil(t, wallet)
}
