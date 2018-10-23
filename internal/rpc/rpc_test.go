package rpc

import (
	"fmt"
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
)

func TestRequest(t *testing.T) {
	//assert := assert.New(t)

	//name := "test"
	params := chaincfg.RegressionNetParams

	fmt.Printf(params.DefaultPort)
	fmt.Printf("testing")

	// _, err := New(name, params, nRpoints)
	//assert.Nil(err)
}
