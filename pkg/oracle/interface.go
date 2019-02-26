package oracle

import (
	"encoding/hex"
	"encoding/json"

	"github.com/btcsuite/btcd/btcec"
)

// PubkeySet contains oracle's pub key and keys for all rate
type PubkeySet struct {
	Pubkey           *btcec.PublicKey
	CommittedRpoints []*btcec.PublicKey
}

// ToJSON serialize PubkeySet to JSON
func (pubset *PubkeySet) ToJSON() ([]byte, error) {
	var rpoints []string
	for _, R := range pubset.CommittedRpoints {
		rpoints = append(rpoints, pubkeyToStr(R))
	}

	v := map[string]interface{}{
		"pubkey":  pubkeyToStr(pubset.Pubkey),
		"rpoints": rpoints,
	}

	s, err := json.Marshal(v)
	return s, err
}

func pubkeyToStr(pub *btcec.PublicKey) string {
	return hex.EncodeToString(pub.SerializeCompressed())
}

// SignSet contains fixed messages and signs
type SignSet struct {
	Msgs  [][]byte
	Signs [][]byte
}
