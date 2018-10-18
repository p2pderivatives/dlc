package oracle

import (
	"fmt"
	"time"

	"github.com/btcsuite/btcd/btcec"
)

// KeySet contains oracle's pub key and keys for all rate
type KeySet struct {
	Pubkey string   `json:"pubkey"`
	Keys   []string `json:"keys"`
}

// KeySet returns a key set for given fixing time
func (oracle Oracle) KeySet(date time.Time) {

}

// deriveKeys derives private and public keys using  hierarchical deterministic format
func (oracle *Oracle) deriveKeys(hdpath ...int) (*btcec.PrivateKey, *btcec.PublicKey, error) {
	var err error

	key := oracle.extKey
	if key == nil {
		err = fmt.Errorf("Extended key must exist")
		return nil, nil, err
	}

	// follow the HD path
	for _, i := range hdpath {
		key, err = key.Child(uint32(i))
		if err != nil {
			return nil, nil, err
		}
	}

	return prvKey, pubKey, nil
}

func ecKeys(key btcec.PrivateKey) (*btcec.PrivateKey, *btcec.PublicKey, error) {
	var prvKey *btcec.PrivateKey
	prvKey, err = key.ECPrivKey()
	if err != nil {
		return nil, nil, err
	}

	var pubKey *btcec.PublicKey
	if pubKey, err = key.ECPubKey(); err != nil {
		return nil, nil, err
	}

}
