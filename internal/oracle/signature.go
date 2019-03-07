package oracle

import (
	"time"

	"github.com/p2pderivatives/dlc/pkg/oracle"
	"github.com/p2pderivatives/dlc/pkg/schnorr"
)

// SignedMsg is an alias of oracle.SignedMsg
type SignedMsg = oracle.SignedMsg

// SignMsg returns FixedMsg for given fixing time
func (oracle *Oracle) SignMsg(ftime time.Time) (SignedMsg, error) {
	msgs, err := oracle.msgsAt(ftime)
	if err != nil {
		return SignedMsg{}, err
	}

	extKey, err := oracle.extKeyForFixingTime(ftime)
	if err != nil {
		return SignedMsg{}, err
	}

	sigs, err := signMsgs(msgs, extKey)
	if err != nil {
		return SignedMsg{}, err
	}

	return SignedMsg{Msgs: msgs, Sigs: sigs}, nil
}

func signMsgs(msgs [][]byte, extKey *privExtKey) ([][]byte, error) {
	opriv, err := extKey.ECPrivKey()
	if err != nil {
		return [][]byte{}, err
	}

	sigs := [][]byte{}

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

		sigs = append(sigs, sign)
	}

	return sigs, nil
}
