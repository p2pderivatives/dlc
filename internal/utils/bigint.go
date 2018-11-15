package utils

import (
	"math/big"

	"github.com/btcsuite/btcd/btcec"
)

// AddBigInts adds signs
func AddBigInts(a, b *big.Int) *big.Int {
	ab := new(big.Int).Add(a, b)
	return new(big.Int).Mod(ab, btcec.S256().N)
}
