package oracle

import (
	"encoding/hex"
	"encoding/json"
	"math"

	"github.com/btcsuite/btcd/btcec"
)

// PubkeySet contains oracle's pub key and keys for all rate
type PubkeySet struct {
	Pubkey           *btcec.PublicKey
	CommittedRpoints []*btcec.PublicKey
}

// PubkeySetJSON is serialized PubkeySet
type PubkeySetJSON struct {
	Pubkey           string   `json:"pubkey"`
	CommittedRpoints []string `json:"rpoints"`
}

// MarshalJSON serialize PubkeySet to JSON
func (pubset PubkeySet) MarshalJSON() ([]byte, error) {
	pubkey := pubkeyToStr(pubset.Pubkey)
	var rpoints []string
	for _, R := range pubset.CommittedRpoints {
		rpoints = append(rpoints, pubkeyToStr(R))
	}

	s, err := json.Marshal(&PubkeySetJSON{
		Pubkey:           pubkey,
		CommittedRpoints: rpoints,
	})
	return s, err
}

func pubkeyToStr(pub *btcec.PublicKey) string {
	return hex.EncodeToString(pub.SerializeCompressed())
}

// UnmarshalJSON deserialize JSON to PubkeySet
func (pubset *PubkeySet) UnmarshalJSON(data []byte) error {
	pjson := &PubkeySetJSON{}
	json.Unmarshal(data, pjson)

	pubkey, err := strToPubkey(pjson.Pubkey)
	if err != nil {
		return err
	}

	var rpoints []*btcec.PublicKey
	for _, rstr := range pjson.CommittedRpoints {
		r, err := strToPubkey(rstr)
		if err != nil {
			return err
		}
		rpoints = append(rpoints, r)
	}

	pubset.Pubkey = pubkey
	pubset.CommittedRpoints = rpoints

	return nil
}

func strToPubkey(str string) (*btcec.PublicKey, error) {
	b, err := hex.DecodeString(str)
	if err != nil {
		return nil, err
	}
	return btcec.ParsePubKey(b, btcec.S256())
}

// SignSet contains fixed messages and signs
type SignSet struct {
	Msgs  [][]byte
	Signs [][]byte
}

// MarshalJSON serialize SignSet to JSON
func (sigset SignSet) MarshalJSON() ([]byte, error) {
	value := ByteMsgsToNumber(sigset.Msgs)

	var sigs []string
	for _, s := range sigset.Signs {
		sigs = append(sigs, hex.EncodeToString(s))
	}

	v := map[string]interface{}{
		"value": value,
		"sigs":  sigs,
	}

	s, err := json.Marshal(v)
	return s, err
}

// NumberToByteMsgs converts number value to byte messages
func NumberToByteMsgs(v int, nDigits int) [][]byte {
	msgs := [][]byte{}

	for i := 0; i < nDigits; i++ {
		d := int(math.Pow(10, float64(nDigits-1-i)))
		b := []byte{byte(v / d)}
		msgs = append(msgs, b)
		v = v % d
	}

	return msgs
}

// ByteMsgsToNumber converts byte messages to number value
func ByteMsgsToNumber(msgs [][]byte) int {
	n := len(msgs)

	v := 0
	for i, m := range msgs {
		d := int(math.Pow(10, float64(n-i-1)))
		v += int(m[0]) * d
	}

	return v
}
