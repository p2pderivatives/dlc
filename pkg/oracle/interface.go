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

// ToJSON serialize SignSet to JSON
func (sigset *SignSet) ToJSON() ([]byte, error) {
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
