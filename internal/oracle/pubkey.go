package oracle

import (
	"encoding/hex"
	"encoding/json"
	"math/big"
	"time"

	"github.com/btcsuite/btcd/btcec"
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

// Commit commits to a message
func (oracle *Oracle) Commit(R *btcec.PublicKey, O *btcec.PublicKey, m []byte) *btcec.PublicKey {
	// H(R,m)
	h := hash(R, m)
	// - H(R,m)
	h = new(big.Int).Mod(new(big.Int).Neg(h), btcec.S256().N)
	hO := new(btcec.PublicKey)
	// - H(R,m)O
	hO.X, hO.Y = btcec.S256().ScalarMult(O.X, O.Y, h.Bytes())
	// R - H(R,m)O
	P := new(btcec.PublicKey)
	P.X, P.Y = btcec.S256().Add(R.X, R.Y, hO.X, hO.Y)
	return P
}

// StrToPubkey converts string to public key.
func StrToPubkey(str string) (*btcec.PublicKey, error) {
	bs, err := hex.DecodeString(str)
	if err != nil {
		return nil, err
	}
	return btcec.ParsePubKey(bs, btcec.S256())
}
