package rpc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBtcdRPC(t *testing.T) {
	rpc, _ := NewBtcdRPC("anyport", "rpcuser", "rpcpassword")
	assert.NotNil(t, rpc)
}
