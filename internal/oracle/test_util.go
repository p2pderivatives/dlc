package oracle

import (
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
)

func newTestOracle(t *testing.T, nRpoints int) (*Oracle, error) {
	name := "test"
	params := chaincfg.RegressionNetParams

	return New(name, params, nRpoints)
}
