package oracle

import (
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
)

func TestNew(t *testing.T) {
	name := "test"
	params := chaincfg.RegressionNetParams
	digit := 1

	_, err := New(name, params, digit)
	if err != nil {
		t.Errorf("Failed to create oracle: %v", err)
		return
	}
}
