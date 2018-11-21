package rpc

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	assert := assert.New(t)

	client, err := NewClient(testConfPath)
	assert.Nil(err)

	_, err = client.RawRequest("getblockchaininfo", []json.RawMessage{})
	assert.NoError(err)
}
