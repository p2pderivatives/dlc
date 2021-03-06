package rpc

import (
	"path/filepath"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcutil"
)

var (
	projectDir, _ = filepath.Abs("../../")
	bitcoinDir    = filepath.Join(projectDir, "bitcoind/")
	testConfName  = "bitcoin.regtest.conf"
	testConfPath  = filepath.Join(bitcoinDir, testConfName)
)

// NewTestRPCClient creates a RPC client for testing
func NewTestRPCClient() (Client, error) {
	return NewClient(testConfPath)
}

// Faucet sends bitcion to a given address
func Faucet(addr btcutil.Address, amt btcutil.Amount) error {
	c, err := NewTestRPCClient()
	if err != nil {
		return err
	}
	_, err = c.SendToAddress(addr, amt)
	if err != nil {
		return err
	}

	_, err = c.Generate(1)
	return err
}

// Generate generates n blocks in regtest
func Generate(n uint32) ([]*chainhash.Hash, error) {
	c, err := NewTestRPCClient()
	if err != nil {
		return nil, err
	}
	return c.Generate(n)
}

// GetBlockCount returns the current block height
func GetBlockCount() (int64, error) {
	c, err := NewTestRPCClient()
	if err != nil {
		return 0, err
	}
	return c.GetBlockCount()
}
