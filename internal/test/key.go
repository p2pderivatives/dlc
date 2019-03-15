package test

import (
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/hdkeychain"
)

// RandKeys generates random key pairs for testing
func RandKeys() (*btcec.PrivateKey, *btcec.PublicKey) {
	seed, _ := hdkeychain.GenerateSeed(hdkeychain.MinSeedBytes)
	extKey, _ := hdkeychain.NewMaster(seed, &chaincfg.RegressionNetParams)
	pub, _ := extKey.ECPubKey()
	priv, _ := extKey.ECPrivKey()
	return priv, pub
}

func RandAddress() *btcutil.AddressWitnessPubKeyHash {
	_, pub := RandKeys()
	pubKeyHash := btcutil.Hash160(pub.SerializeCompressed())
	addr, _ := btcutil.NewAddressWitnessPubKeyHash(
		pubKeyHash, &chaincfg.RegressionNetParams)
	return addr
}
