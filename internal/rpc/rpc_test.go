package rpc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testport = "localhost:18433" //18443 for regnet, 18332 for testnet3
	testuser = "username"
	testpass = "password"
)

func TestNewClient(t *testing.T) {
	client, err := NewClient(testport, testuser, testpass)
	//defer client.Shutdown()

	assert.Nil(t, err)
	assert.NotNil(t, client)
}
