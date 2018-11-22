package utils

import "github.com/btcsuite/btcutil"

// ItoAmt converts int to btcutil.Amount
func ItoAmt(n int) btcutil.Amount {
	return btcutil.Amount(float64(n))
}
