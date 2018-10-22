package oracle

import (
	"encoding/hex"
	"encoding/json"
	"math/big"
	"time"

	"github.com/dgarage/dlc/internal/schnorr"
)

// SignSet contains fixing value, messages and sign for each committed R point
type SignSet struct {
	Value string   `json:"value"`
	Msgs  []string `json:"msgs"`
	Signs []string `json:"signs"`
}

// ToJSON encodes signset to json
func (k *SignSet) ToJSON() ([]byte, error) {
	return json.Marshal(k)
}

// SignSet returns SignSet for given fixing time
func (oracle *Oracle) SignSet(ftime time.Time) (SignSet, error) {
	vals, err := oracle.valuesAt(ftime)
	if err != nil {
		return SignSet{}, err
	}

	extKey, err := oracle.extKeyForFixingTime(ftime)
	if err != nil {
		return SignSet{}, err
	}

	msgs, signs, err := signToValues(vals, extKey)
	if err != nil {
		return SignSet{}, err
	}

	return SignSet{"", msgs, signs}, nil

}

func signToValues(vals []int, extKey *privExtKey) ([]string, []string, error) {
	opriv, err := extKey.ECPrivKey()
	if err != nil {
		return []string{}, []string{}, err
	}

	msgs := []string{}
	signs := []string{}
	for i, val := range vals {
		k, err := extKey.derive(i)
		if err != nil {
			return []string{}, []string{}, err
		}
		rpriv, err := k.ECPrivKey()
		if err != nil {
			return []string{}, []string{}, err
		}

		m := big.NewInt(int64(val)).Bytes()

		// Schnorr signature
		s := schnorr.Sign(rpriv, opriv, m)

		signs = append(signs, hex.EncodeToString(s))
		msgs = append(msgs, hex.EncodeToString(m))
	}

	return msgs, signs, nil
}
