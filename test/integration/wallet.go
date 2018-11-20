package integration

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
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
	walletDir, err := ioutil.TempDir("", randomWalletDirPrefix(seed))
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

func randomWalletDirPrefix(seed []byte) string {
	hashB := md5.Sum(seed)
	hash := hex.EncodeToString(hashB[:])
	return fmt.Sprintf("dlcwallet_%s", hash)
}
