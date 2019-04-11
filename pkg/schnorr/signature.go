package schnorr

import (
	"crypto/sha256"
	"math/big"

	"github.com/btcsuite/btcd/btcec"
)

// CommitMulti calculates a commitment by summing commitments of multiple msgs
func CommitMulti(
	V *btcec.PublicKey, Rs []*btcec.PublicKey, msgs [][]byte,
) *btcec.PublicKey {

	Psum := new(btcec.PublicKey)
	for i, m := range msgs {
		R := Rs[i]
		P := Commit(V, R, m)
		Psum = addPubkeys(Psum, P)
	}
	return Psum
}

// Commit is calculatd by the following formula
//   sG = R - h(R, m) * V
// Where
//   s: sign for the message m
//   G: elliptic curve base
//   R: R-point
//   m: message
//   V: oracle's public key
func Commit(V, R *btcec.PublicKey, m []byte) *btcec.PublicKey {
	// - h(R, m)
	h := hash(R, m)
	h = new(big.Int).Neg(h)
	h = new(big.Int).Mod(h, btcec.S256().N)

	// - h(R, m) * V
	hV := new(btcec.PublicKey)
	hV.X, hV.Y = btcec.S256().ScalarMult(V.X, V.Y, h.Bytes())

	// R - h(R, m) * V
	P := addPubkeys(R, hV)
	return P
}

func addPubkeys(A, B *btcec.PublicKey) *btcec.PublicKey {
	C := new(btcec.PublicKey)
	if A.X == nil {
		C.X, C.Y = B.X, B.Y
	} else if B.X == nil {
		C.X, C.Y = A.X, A.Y
	} else {
		C.X, C.Y = btcec.S256().Add(A.X, A.Y, B.X, B.Y)
	}
	return C
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
//   rpriv: random point EC private key opriv: oracle's EC private key
//   m: message
func Sign(opriv, rpriv *btcec.PrivateKey, m []byte) []byte {
	R := rpriv.PubKey()
	k := rpriv.D
	v := opriv.D

	// h(R,m) * v
	hv := new(big.Int).Mul(hash(R, m), v)

	// k - h(R,m) * v
	s := new(big.Int).Sub(k, hv)

	// s mod N
	s = new(big.Int).Mod(s, btcec.S256().N)

	return s.Bytes()
}

// SumSigs sums signs up for a multi-message commitment
func SumSigs(signs [][]byte) []byte {
	sum := new(big.Int)
	for _, sign := range signs {
		sb := new(big.Int).SetBytes(sign)
		sum = new(big.Int).Add(sum, sb)
	}
	sum = new(big.Int).Mod(sum, btcec.S256().N)
	return sum.Bytes()
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
func Verify(P *btcec.PublicKey, sign []byte) bool {
	sG := new(btcec.PublicKey)
	sG.X, sG.Y = btcec.S256().ScalarBaseMult(sign)
	return P.IsEqual(sG)
}
