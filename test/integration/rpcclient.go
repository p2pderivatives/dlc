package integration

import (
	"path/filepath"

	"github.com/dgarage/dlc/internal/rpc"
)

var (
	projectDir, _ = filepath.Abs("../../")
	bitcoinDir    = filepath.Join(projectDir, "bitcoind/")
	btcconfName   = "bitcoin.regtest.conf"
	btcconfPath   = filepath.Join(bitcoinDir, btcconfName)
)

// NewRPCClient creates rpcclient for integration testing
func NewRPCClient() (rpc.Client, error) {
	return rpc.NewClient(btcconfPath)
}
