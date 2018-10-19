package oracle

import (
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/btcsuite/btcutil/hdkeychain"
)

// KeySet contains oracle's pub key and keys for all rate
type KeySet struct {
	Pubkey           string   `json:"pubkey"`
	CommittedRpoints []string `json:"committed_rpoints"`
}

// ToJSON encodes keyset to json
func (k KeySet) ToJSON() ([]byte, error) {
	return json.Marshal(k)
}

// Extended key wrapper
type privExtKey struct {
	key *hdkeychain.ExtendedKey
}

func (key privExtKey) pubKeyStr() (string, error) {
	pubkey, err := key.key.ECPubKey()
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(pubkey.SerializeCompressed()), nil
}

func (oracle Oracle) baseKey() privExtKey {
	// TODO: define HD path following bip44, 47
	return privExtKey{oracle.masterKey}
}

// KeySet returns a key set for given fixing time
// TODO: Add a document for keyset generation
func (oracle Oracle) KeySet(ftime time.Time) (KeySet, error) {
	// TODO: Should we check if it's later than now?

	// derive oracle's pubkey for the given time
	hdpath := timeToHDpath(ftime)
	extKey, err := deriveKeys(oracle.baseKey(), hdpath...)
	if err != nil {
		return KeySet{}, err
	}
	pubkey, err := extKey.pubKeyStr()
	if err != nil {
		return KeySet{}, err
	}

	// derive pubkeys for all committed R-points at the given time
	rpoints, err := committedRpoints(extKey, oracle.nRpoints)
	if err != nil {
		return KeySet{}, err
	}

	keyset := KeySet{pubkey, rpoints}

	return keyset, nil
}

func committedRpoints(extKey privExtKey, nRpoints int) ([]string, error) {
	keys := []string{}
	for i := 0; i < nRpoints; i++ {
		k, err := deriveKeys(extKey, i)
		if err != nil {
			return nil, err
		}
		pubKeyStr, err := k.pubKeyStr()
		if err != nil {
			return nil, err
		}
		keys = append(keys, pubKeyStr)
	}

	return keys, nil
}

func timeToHDpath(t time.Time) []int {
	return []int{t.Year(), int(t.Month()), t.Day(), t.Hour(), t.Minute(), t.Second()}
}

// deriveKeys derives private and public keys using hierarchical deterministic format
func deriveKeys(extKey privExtKey, path ...int) (privExtKey, error) {
	for _, i := range path {
		key, err := extKey.key.Child(uint32(i))
		if err != nil {
			return privExtKey{}, err
		}
		extKey = privExtKey{key}
	}

	return extKey, nil
}
