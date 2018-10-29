package oracle

import (
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/btcsuite/btcd/btcec"
	"github.com/dgarage/dlc/internal/schnorr"
	"github.com/stretchr/testify/assert"
)

func TestSignSetNoMsgsFailure(t *testing.T) {
	assert := assert.New(t)

	o := NewTestOracle()
	ftime := time.Now()

	_, err := o.SignSet(ftime)
	assert.NotNil(err)
}

func TestSignSet(t *testing.T) {
	assert := assert.New(t)

	o := NewTestOracle()
	ftime := time.Now()
	pub, _ := o.PubkeySet(ftime)

	// Fix msgs
	msgs := randomMsgs(o.nRpoints)
	err := o.fixMsgs(ftime, msgs)
	assert.Nil(err)

	// Get signetures
	signSet, err := o.SignSet(ftime)
	assert.Nil(err)
	assert.Equal(len(signSet.Msgs), len(signSet.Signs))
	assert.Equal(len(signSet.Msgs), len(pub.CommittedRpoints))

	// Verification for each commit/sign pair
	V := pub.Pubkey
	Rs := pub.CommittedRpoints
	Ps := []*btcec.PublicKey{}
	for i := 0; i < len(signSet.Msgs); i++ {
		R := Rs[i]
		m := signSet.Msgs[i]
		sign := signSet.Signs[i]
		P := schnorr.Commit(V, R, m)
		Ps = append(Ps, P)

		assert.True(schnorr.Verify(P, sign))
	}

	// Verification for summmed commit/sign pair
	Psum := sumPubkeys(Ps)
	sigsum := sumSigns(signSet.Signs)
	assert.True(schnorr.Verify(Psum, sigsum))
}

func randomMsgs(n int) [][]byte {
	var msgs [][]byte
	for i := 0; i < n; i++ {
		m := big.NewInt(rand.Int63n(9)).Bytes()
		msgs = append(msgs, m)
	}
	return msgs
}

func sumPubkeys(pubs []*btcec.PublicKey) *btcec.PublicKey {
	sum := new(btcec.PublicKey)
	for _, pub := range pubs {
		if sum.X == nil {
			sum.X, sum.Y = pub.X, pub.Y
		} else {
			sum.X, sum.Y = btcec.S256().Add(sum.X, sum.Y, pub.X, pub.Y)
		}
	}
	return sum
}

func sumSigns(signs [][]byte) []byte {
	sum := new(big.Int)
	for _, sign := range signs {
		sb := new(big.Int).SetBytes(sign)
		sum = new(big.Int).Add(sum, sb)
	}
	return sum.Bytes()
}
