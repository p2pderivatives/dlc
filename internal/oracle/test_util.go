package oracle

import "github.com/btcsuite/btcd/chaincfg"

// NewTestOracle creates a oracle for test
func NewTestOracle() *Oracle {
	name := "test"
	params := chaincfg.RegressionNetParams
	nRpoints := 3

	o, _ := New(name, params, nRpoints)
	o.InitDB()

	return o
}
