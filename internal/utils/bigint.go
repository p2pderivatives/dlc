package utils

import (
	"math/big"

	"github.com/btcsuite/btcd/btcec"
)

// AddInts adds signs
func AddInts(a, b *big.Int) *big.Int {
	ab := new(big.Int).Add(a, b)
	return new(big.Int).Mod(ab, btcec.S256().N)
}
