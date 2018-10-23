package schnorr

import (
	"crypto/sha256"
	"math/big"

	"github.com/btcsuite/btcd/btcec"
)

// Commit is calculatd by the following formula
//   sG = R - h(R, m) * V
// Where
//   s: sign for the message m
//   G: elliptic curve base
//   R: R-point
//   m: message
//   V: oracle's public key
func Commit(
	V *btcec.PublicKey, R *btcec.PublicKey, m []byte) *btcec.PublicKey {
	// - h(R, m)
	h := hash(R, m)
	h = new(big.Int).Neg(h)
	h = new(big.Int).Mod(h, btcec.S256().N)

	// - h(R, m) * V
	hV := new(btcec.PublicKey)
	hV.X, hV.Y = btcec.S256().ScalarMult(V.X, V.Y, h.Bytes())

	// R - h(R, m) * V
	P := new(btcec.PublicKey)
	P.X, P.Y = btcec.S256().Add(R.X, R.Y, hV.X, hV.Y)

	return P
}

// Sign is calculated by the following formula
//   s = k - h(R, m) * v
// Where
//   s: sign
//   h: hash function
//   k: random nonce
//   R: R-point R = kG
//   m: message
//   G: elliptic curve base
//   v: oracle's private key
// Parameters:
//   rpriv: random point EC private key
//   opriv: oracle's EC private key
//   m: message
func Sign(
	opriv *btcec.PrivateKey, rpriv *btcec.PrivateKey, m []byte) *big.Int {
	R := rpriv.PubKey()
	k := rpriv.D
	v := opriv.D

	// h(R,m) * v
	hv := new(big.Int).Mul(hash(R, m), v)

	// k - h(R,m) * v
	s := new(big.Int).Sub(k, hv)

	// s mod N
	s = new(big.Int).Mod(s, btcec.S256().N)
	return s
}

func hash(R *btcec.PublicKey, m []byte) *big.Int {
	s := sha256.New()
	s.Write(R.SerializeUncompressed())
	s.Write(m)
	h := new(big.Int).SetBytes(s.Sum(nil))
	h = new(big.Int).Mod(h, btcec.S256().N)
	return h
}

// Verify verfies sG = R - h(R, m) * V
func Verify(P *btcec.PublicKey, sign *big.Int) bool {
	sG := new(btcec.PublicKey)
	sG.X, sG.Y = btcec.S256().ScalarBaseMult(sign.Bytes())
	return P.IsEqual(sG)
}
