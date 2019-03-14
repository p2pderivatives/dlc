package oracle

import (
	"encoding/hex"
	"encoding/json"
	"math"

	"github.com/btcsuite/btcd/btcec"
	"github.com/p2pderivatives/dlc/pkg/utils"
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
	s, err := json.Marshal(pubset.JSON())
	return s, err
}

// JSON returns PubkeySetJSON
func (pubset PubkeySet) JSON() *PubkeySetJSON {
	pubkey := utils.PubkeyToStr(pubset.Pubkey)
	var rpoints []string
	for _, R := range pubset.CommittedRpoints {
		rpoints = append(rpoints, utils.PubkeyToStr(R))
	}
	return &PubkeySetJSON{
		Pubkey:           pubkey,
		CommittedRpoints: rpoints,
	}
}

// UnmarshalJSON deserialize JSON to PubkeySet
func (pubset *PubkeySet) UnmarshalJSON(data []byte) error {
	pjson := &PubkeySetJSON{}
	err := json.Unmarshal(data, pjson)
	if err != nil {
		return err
	}
	return pubset.ParseJSON(pjson)
}

// ParseJSON parses PubkeySetJSON
func (pubset *PubkeySet) ParseJSON(pjson *PubkeySetJSON) error {
	pubkey, err := utils.ParsePublicKey(pjson.Pubkey)
	if err != nil {
		return err
	}

	var rpoints []*btcec.PublicKey
	for _, rstr := range pjson.CommittedRpoints {
		r, err := utils.ParsePublicKey(rstr)
		if err != nil {
			return err
		}
		rpoints = append(rpoints, r)
	}

	pubset.Pubkey = pubkey
	pubset.CommittedRpoints = rpoints

	return nil
}

// SignedMsg contains fixed messages and signatures
type SignedMsg struct {
	Msgs [][]byte
	Sigs [][]byte
}

// SignedMsgJSON is siged message in JSON format
type SignedMsgJSON struct {
	Value int      `json:"value"`
	Sigs  []string `json:"sigs"`
}

// MarshalJSON serialize SignSet to JSON
func (sm SignedMsg) MarshalJSON() ([]byte, error) {
	value := ByteMsgsToNumber(sm.Msgs)

	var sigs []string
	for _, s := range sm.Sigs {
		sigs = append(sigs, hex.EncodeToString(s))
	}

	return json.Marshal(&SignedMsgJSON{
		Value: value,
		Sigs:  sigs,
	})
}

// UnmarshalJSON deserialize JSON to SignedMsg
func (sm *SignedMsg) UnmarshalJSON(data []byte) error {
	smJSON := &SignedMsgJSON{}
	err := json.Unmarshal(data, smJSON)
	if err != nil {
		return err
	}

	n := len(smJSON.Sigs)
	sm.Msgs = NumberToByteMsgs(smJSON.Value, n)

	var sigs [][]byte
	for _, s := range smJSON.Sigs {
		sig, err := hex.DecodeString(s)
		if err != nil {
			return err
		}
		sigs = append(sigs, sig)
	}
	sm.Sigs = sigs

	return nil
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
