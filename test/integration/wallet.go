package integration

import (
	"io/ioutil"
	"path/filepath"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/dgarage/dlc/internal/rpc"
	"github.com/dgarage/dlc/internal/wallet"
)

var (
	projectDir, _ = filepath.Abs("../../")
	bitcoinDir    = filepath.Join(projectDir, "bitcoind/")
	btcconfName   = "bitcoin.regtest.conf"
	btcconfPath   = filepath.Join(bitcoinDir, btcconfName)
)

func newWallet(name string, pubpass, privpass []byte) (wallet.Wallet, error) {
	params := &chaincfg.RegressionNetParams

	// generate random seed
	seed, err := hdkeychain.GenerateSeed(
		hdkeychain.RecommendedSeedLen)
	if err != nil {
		return nil, err
	}

	// create wallet dbdir
	walletDir, err := ioutil.TempDir("", "dlcwallet")
	if err != nil {
		return nil, err
	}

	// create wallet
	w, err := wallet.CreateWallet(
		params, seed, pubpass, privpass, walletDir, name)
	if err != nil {
		return nil, err
	}

	// create rpcclient
	rpcclient, err := rpc.NewClient(btcconfPath)
	if err != nil {
		return nil, err
	}

	w.SetRPCClient(rpcclient)

	return w, nil
}
