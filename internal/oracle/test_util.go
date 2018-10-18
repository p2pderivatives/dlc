package oracle

import (
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
)

func newTestOracle(t *testing.T) (*Oracle, error) {
	name := "test"
	params := chaincfg.RegressionNetParams
	digit := 1

	return New(name, params, digit)
}
