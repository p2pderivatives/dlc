package oracle

import (
	"encoding/hex"
	"time"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil/hdkeychain"
)

// Extended key wrapper
type privExtKey struct {
	key *hdkeychain.ExtendedKey
}

func (oracle *Oracle) baseKey() privExtKey {
	// TODO: define HD path following bip44, 47
	return privExtKey{oracle.masterKey}
}

func (key *privExtKey) ECPubKey() (*btcec.PublicKey, error) {
	return key.key.ECPubKey()
}

func (key *privExtKey) ECPrivKey() (*btcec.PrivateKey, error) {
	return key.key.ECPrivKey()
}

func (key *privExtKey) pubKeyStr() (string, error) {
	pubkey, err := key.ECPubKey()
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(pubkey.SerializeCompressed()), nil
}

// deriveKeys derives child key following HD path
func (key privExtKey) derive(path ...int) (*privExtKey, error) {
	for _, i := range path {
		extKey, err := key.key.Child(uint32(i))
		if err != nil {
			return nil, err
		}
		key = privExtKey{extKey}
	}

	return &key, nil
}

func (oracle *Oracle) extKeyForFixingTime(ftime time.Time) (*privExtKey, error) {
	hdpath := timeToHDpath(ftime)
	baseKey := oracle.baseKey()
	return baseKey.derive(hdpath...)
}

func timeToHDpath(t time.Time) []int {
	return []int{t.Year(), int(t.Month()), t.Day(), t.Hour(), t.Minute(), t.Second()}
}
