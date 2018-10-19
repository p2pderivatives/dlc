package wallet

import (
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
)

const (
	// RecommendedSeedLen is the recommended length in bytes for a seed
	// to a master node.
	RecommendedSeedLen = 32 // 256 bits
)

func TestNewWallet(t *testing.T) {
	params := chaincfg.RegressionNetParams
	seed, _ := hdkeychain.GenerateSeed(RecommendedSeedLen)

	_, err := NewWallet(params, seed)
	if err != nil {
		t.Errorf("Failed to create wallet: %v", err)
		return
	}
}
