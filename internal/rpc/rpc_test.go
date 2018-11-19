package rpc

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	projectDir, _ = filepath.Abs("../../")
	bitcoinDir    = filepath.Join(projectDir, "bitcoind/")
	confName      = "bitcoin.regtest.conf"
	confPath      = filepath.Join(bitcoinDir, confName)
)

func TestNewClient(t *testing.T) {
	assert := assert.New(t)

	client, err := NewClient(confPath)
	assert.Nil(err)

	_, err = client.ListUnspent()
	assert.NoError(err)
}
