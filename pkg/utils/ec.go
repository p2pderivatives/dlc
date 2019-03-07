package utils

import (
	"encoding/hex"

	"github.com/btcsuite/btcd/btcec"
)

// PubkeyToStr converts btcec.PublicKey to hex string (compressed)
func PubkeyToStr(pub *btcec.PublicKey) string {
	return hex.EncodeToString(pub.SerializeCompressed())
}

// ParsePublicKey converts hex string to btcec.PublicKey
func ParsePublicKey(pub string) (*btcec.PublicKey, error) {
	b, err := hex.DecodeString(pub)
	if err != nil {
		return nil, err
	}
	return btcec.ParsePubKey(b, btcec.S256())
}
