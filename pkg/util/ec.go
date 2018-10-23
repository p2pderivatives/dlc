package util

import (
	"encoding/hex"

	"github.com/btcsuite/btcd/btcec"
)

// StirngToECPubkey deserializes hex string to EC public key
func StirngToECPubkey(str string) (*btcec.PublicKey, error) {
	bs, err := hex.DecodeString(str)
	if err != nil {
		return nil, err
	}
	return btcec.ParsePubKey(bs, btcec.S256())
}

// ECPubKeyToString serialized EC publick key to hex string
func ECPubKeyToString(pubkey *btcec.PublicKey) string {
	return hex.EncodeToString(pubkey.SerializeCompressed())
}
