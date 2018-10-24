package rpc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	// RecommendedSeedLen is the recommended length in bytes for a seed
	// to a master node.
	RecommendedSeedLen = 32 // 256 bits

)

func TestNewBtcRPC(t *testing.T) {
	rpc := NewBtcRPC("anyport", "rpcuser", "rpcpassword")
	assert.NotNil(t, rpc)
}

// TODO
func TestRequest(t *testing.T) {
	// mock bitcoind/testnet or testnet?
	//testy := NewBtcRPC("http://localhost:18443", "user", "pass")
	//test_rpc := NewBtcRPC("http://localhost:18332", "akek", "akek")
	//res2, _ := test_rpc.Request("getblockchaininfo")

	//fmt.Printf("%+v\n", res)
	//fmt.Printf("%+v\n", res2)
}
