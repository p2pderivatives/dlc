package schnorr

import (
	"math/big"
	"testing"

	"github.com/btcsuite/btcd/btcec"
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

	// verifiation
	c := btcec.S256().CurveParams
	X, Y := btcec.S256().ScalarMult(c.Gx, c.Gy, sign.Bytes())
	assert.Equal(sG.X, X)
	assert.Equal(sG.Y, Y)

	// Oracle's sign for another message
	m = big.NewInt(int64(2)).Bytes()
	sign = Sign(opriv, rpriv, m)
	X, Y = btcec.S256().ScalarMult(c.Gx, c.Gy, sign.Bytes())
	assert.NotEqual(sG.X, X)
	assert.NotEqual(sG.Y, Y)
}

func randExtKey() (*hdkeychain.ExtendedKey, error) {
	seed, err := hdkeychain.GenerateSeed(hdkeychain.MinSeedBytes)
	if err != nil {
		return nil, err
	}
	return hdkeychain.NewMaster(seed, &chaincfg.RegressionNetParams)
}
