package test

import (
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
)

func RandKeys() (*btcec.PrivateKey, *btcec.PublicKey) {
	seed, _ := hdkeychain.GenerateSeed(hdkeychain.MinSeedBytes)
	extKey, _ := hdkeychain.NewMaster(seed, &chaincfg.RegressionNetParams)
	pub, _ := extKey.ECPubKey()
	priv, _ := extKey.ECPrivKey()
	return priv, pub
}
