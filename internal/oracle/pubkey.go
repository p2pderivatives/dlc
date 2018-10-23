package oracle

import (
	"encoding/json"
	"time"
)

// PubkeySet contains oracle's pub key and keys for all rate
type PubkeySet struct {
	Pubkey           string   `json:"pubkey"`
	CommittedRpoints []string `json:"committed_rpoints"`
}

// ToJSON encodes pubkey set to json
func (k *PubkeySet) ToJSON() ([]byte, error) {
	return json.Marshal(k)
}

// PubkeySet returns a key set for given fixing time
// TODO: Add a document for pubkey set generation
func (oracle *Oracle) PubkeySet(ftime time.Time) (PubkeySet, error) {
	// TODO: Should we check if it's later than now?

	extKey, err := oracle.extKeyForFixingTime(ftime)
	// derive oracle's pubkey for the given time
	if err != nil {
		return PubkeySet{}, err
	}
	pubkey, err := extKey.pubKeyStr()
	if err != nil {
		return PubkeySet{}, err
	}

	// derive pubkeys for all committed R-points at the given time
	rpoints, err := committedRpoints(extKey, oracle.nRpoints)
	if err != nil {
		return PubkeySet{}, err
	}

	keyset := PubkeySet{pubkey, rpoints}

	return keyset, nil
}

func committedRpoints(extKey *privExtKey, nRpoints int) ([]string, error) {
	keys := []string{}
	for i := 0; i < nRpoints; i++ {
		k, err := extKey.derive(i)
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
