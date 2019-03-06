package oracle

import (
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/dgarage/dlc/pkg/schnorr"
	"github.com/stretchr/testify/assert"
)

func TestSignSetNoMsgsFailure(t *testing.T) {
	assert := assert.New(t)

	o := NewTestOracle()
	ftime := time.Now()

	_, err := o.SignMsg(ftime)
	assert.NotNil(err)
}

func TestSignSet(t *testing.T) {
	assert := assert.New(t)

	o := NewTestOracle()
	ftime := time.Now()
	pub, _ := o.PubkeySet(ftime)

	// Fix msgs
	msgs := randomMsgs(o.nRpoints)
	err := o.FixMsgs(ftime, msgs)
	assert.Nil(err)

	// Get signetures
	signSet, err := o.SignMsg(ftime)
	assert.Nil(err)
	assert.Equal(len(signSet.Msgs), len(signSet.Sigs))
	assert.Equal(len(signSet.Msgs), len(pub.CommittedRpoints))

	// Verification for each commit/sign pair
	V := pub.Pubkey
	Rs := pub.CommittedRpoints
	for i := 0; i < len(signSet.Msgs); i++ {
		R := Rs[i]
		m := signSet.Msgs[i]
		sign := signSet.Sigs[i]
		P := schnorr.Commit(V, R, m)
		assert.True(schnorr.Verify(P, sign))
	}

	// Verification for summmed commit/sign pair
	Psum := schnorr.CommitMulti(V, Rs, signSet.Msgs)
	sigsum := schnorr.SumSigs(signSet.Sigs)
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
