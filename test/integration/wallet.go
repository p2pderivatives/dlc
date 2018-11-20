package integration

import (
	"io/ioutil"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/dgarage/dlc/internal/wallet"
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
	rpcclient, err := NewRPCClient()
	if err != nil {
		return nil, err
	}

	w.SetRPCClient(rpcclient)

	return w, nil
}
