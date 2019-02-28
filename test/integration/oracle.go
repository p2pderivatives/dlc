package integration

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/dgarage/dlc/internal/oracle"
)

// NewOracle creates an oracle for integration tests
func newOracle(name string, nPoints int) (*oracle.Oracle, error) {
	params := &chaincfg.RegressionNetParams

	o, err := oracle.New(name, params, nPoints)
	if err != nil {
		return nil, err
	}
	o.InitDB()
	return o, nil
}
