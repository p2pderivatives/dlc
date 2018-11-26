package schnorr

import (
	"math/big"
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/stretchr/testify/assert"
)

func TestSchnorrSignature(t *testing.T) {
	assert := assert.New(t)

	// Oracle's private keys
	extKey, _ := randExtKey()
	extKeyChild, _ := extKey.Child(1)

	// Public keys
	V, _ := extKey.ECPubKey()
	R, _ := extKeyChild.ECPubKey()

	// message
	m := big.NewInt(int64(1)).Bytes()
	sG := Commit(V, R, m)

	// Oracle's sign for the committed message
	opriv, _ := extKey.ECPrivKey()
	rpriv, _ := extKeyChild.ECPrivKey()
	sign := Sign(opriv, rpriv, m)

	// should pass verifiation
	assert.True(Verify(sG, sign))

	// Oracle's sign for another message
	m2 := big.NewInt(int64(2)).Bytes()
	sign2 := Sign(opriv, rpriv, m2)

	// should not pass verifiation
	assert.False(Verify(sG, sign2))
}

func randExtKey() (*hdkeychain.ExtendedKey, error) {
	seed, err := hdkeychain.GenerateSeed(hdkeychain.MinSeedBytes)
	if err != nil {
		return nil, err
	}
	return hdkeychain.NewMaster(seed, &chaincfg.RegressionNetParams)
}
