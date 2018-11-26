package oracle

import (
	"time"

	"github.com/dgarage/dlc/pkg/oracle"
	"github.com/dgarage/dlc/pkg/schnorr"
)

type SignSet = oracle.SignSet

// SignSet returns SignSet for given fixing time
func (oracle *Oracle) SignSet(ftime time.Time) (SignSet, error) {
	msgs, err := oracle.msgsAt(ftime)
	if err != nil {
		return SignSet{}, err
	}

	extKey, err := oracle.extKeyForFixingTime(ftime)
	if err != nil {
		return SignSet{}, err
	}

	signs, err := signMsgs(msgs, extKey)
	if err != nil {
		return SignSet{}, err
	}

	return SignSet{msgs, signs}, nil
}

func signMsgs(msgs [][]byte, extKey *privExtKey) ([][]byte, error) {
	opriv, err := extKey.ECPrivKey()
	if err != nil {
		return [][]byte{}, err
	}

	signs := [][]byte{}

	for i, m := range msgs {
		k, err := extKey.derive(i)
		if err != nil {
			return [][]byte{}, err
		}
		rpriv, err := k.ECPrivKey()
		if err != nil {
			return [][]byte{}, err
		}

		// Schnorr signature
		sign := schnorr.Sign(opriv, rpriv, m)

		signs = append(signs, sign)
	}

	return signs, nil
}
