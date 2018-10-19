package oracle

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil/hdkeychain"
)

// KeySet contains oracle's pub key and keys for all rate
type KeySet struct {
	Pubkey       string   `json:"pubkey"`
	DigitPubkeys []string `json:"digit_pubkeys"`
}

// ToJSON encodes keyset to json
func (k KeySet) ToJSON() ([]byte, error) {
	return json.Marshal(k)
}

// KeySet returns a key set for given fixing time
func (oracle Oracle) KeySet(ftime time.Time) (KeySet, error) {
	// TODO: Should we check if it's later than now?

	// derive oracle's pubkey for the given time
	hdpath := timeToHDpath(ftime)
	_, pubkey, err := oracle.deriveKeys(hdpath...)
	if err != nil {
		return KeySet{}, err
	}

	// derive pubkeys for all digits at the given time
	dPubkeys, err := oracle.digitPubkeys(hdpath)
	if err != nil {
		return KeySet{}, err
	}

	keyset := KeySet{pubKeyToString(pubkey), dPubkeys}

	return keyset, nil
}

func (oracle Oracle) digitPubkeys(hdpath []int) ([]string, error) {
	keys := []string{}
	for i := 0; i < oracle.digit; i++ {
		_, key, err := oracle.deriveKeys(append(hdpath, i)...)
		if err != nil {
			return nil, err
		}
		keys = append(keys, pubKeyToString(key))
	}

	return keys, nil
}

func timeToHDpath(t time.Time) []int {
	return []int{t.Year(), int(t.Month()), t.Day(), t.Hour(), t.Minute(), t.Second()}
}

func pubKeyToString(key *btcec.PublicKey) string {
	return hex.EncodeToString(key.SerializeCompressed())
}

// deriveKeys derives private and public keys using hierarchical deterministic format
func (oracle *Oracle) deriveKeys(hdpath ...int) (
	*btcec.PrivateKey, *btcec.PublicKey, error,
) {
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

	return ecKeys(key)
}

func ecKeys(key *hdkeychain.ExtendedKey) (
	*btcec.PrivateKey, *btcec.PublicKey, error,
) {
	prvKey, err := key.ECPrivKey()
	if err != nil {
		return nil, nil, err
	}

	pubKey, err := key.ECPubKey()
	if err != nil {
		return nil, nil, err
	}

	return prvKey, pubKey, nil
}
