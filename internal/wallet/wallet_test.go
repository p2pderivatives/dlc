package wallet

import (
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/stretchr/testify/assert"
)

const (
	// RecommendedSeedLen is the recommended length in bytes for a seed
	// to a master node.
	RecommendedSeedLen = 32 // 256 bits
)

func TestNewWallet(t *testing.T) {
	params := chaincfg.RegressionNetParams
	seed, _ := hdkeychain.GenerateSeed(RecommendedSeedLen)

	wallet, _ := NewWallet(params, seed)
	assert.NotNil(t, wallet)
}

func TestGetWitnessSignature(t *testing.T) {
}
