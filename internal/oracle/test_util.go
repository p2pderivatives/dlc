package oracle

import (
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
)

func newTestOracle(t *testing.T) (*Oracle, error) {
	name := "test"
	params := chaincfg.RegressionNetParams
	nRpoints := 1

	return New(name, params, digit)
}
