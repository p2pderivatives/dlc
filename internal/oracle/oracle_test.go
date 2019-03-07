package oracle

import (
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert := assert.New(t)

	name := "test"
	params := &chaincfg.RegressionNetParams
	nRpoints := 1

	_, err := New(name, params, nRpoints)
	assert.Nil(err)
}
