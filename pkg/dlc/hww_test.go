package dlc

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/p2pderivatives/dlc/internal/test"
	"github.com/p2pderivatives/dlc/pkg/script"
	"github.com/stretchr/testify/assert"
)

const SigHashType = txscript.SigHashAll

func TestHWW(t *testing.T) {
	assert := assert.New(t)

	pubstr := "xpub6HAvEkMsGKfWdrFB7ErzZX5YLFWMwJUZpaX1VWEzbV9VLPRgBfNqaJzZWgZTCm9jnqKEzCjMVJ7g1REoYPetqbWFzAgyVsFpjUpWnkVuJmm"
	key, err := hdkeychain.NewKeyFromString(pubstr)
	assert.NoError(err)

	pub, err := key.ECPubKey()
	assert.NoError(err)

	// oracle's commitment and signature
	oprivHex := "c53260a779e799341271547b20e0092974a9141bfba6cc574dd10654d5775524"
	oprivB, err := hex.DecodeString(oprivHex)
	assert.NoError(err)
	_, opub := btcec.PrivKeyFromBytes(btcec.S256(), oprivB)

	// prepare source/redeem tx
	amt := int64(10000)
	sourceTx := test.NewSourceTx()
	addedPub := addPubkeys(pub, opub)
	sc, err := script.P2WPKHpkScript(addedPub)
	assert.NoError(err)
	sourceTx.AddTxOut(wire.NewTxOut(amt, sc))
	redeemTx := test.NewRedeemTx(sourceTx, 0)

	sigHash := txscript.NewTxSigHashes(redeemTx)
	wsHash, err := txscript.CalcWitnessSigHash(
		sc, sigHash, SigHashType, redeemTx, 0, amt)
	assert.NoError(err)
	// fmt.Printf("Witness SigHash: %x\n", wsHash)

	sigstr := "c5245997c396954406b22437f7785cfb145c8aa807ecb302ca2f3ac3c54b7b1b14edc29bb779cdd452b42dd65cb21cca39eab39451e09f06abd2027cf2426216"
	sigB, err := hex.DecodeString(sigstr)
	assert.NoError(err)
	sig := decodeRawSignature(sigB)

	verified := sig.Verify(wsHash, addedPub)
	assert.True(verified)

	witsig := append(sig.Serialize(), byte(SigHashType))
	witness := wire.TxWitness{witsig, addedPub.SerializeCompressed()}
	redeemTx.TxIn[0].Witness = witness

	err = test.ExecuteScript(sc, redeemTx, amt)
	assert.NoError(err)
}

func addPubkeys(A, B *btcec.PublicKey) *btcec.PublicKey {
	if A.X == nil {
		return B
	} else if B.X == nil {
		return A
	} else {
		A.X, A.Y = btcec.S256().Add(A.X, A.Y, B.X, B.Y)
		return A
	}
}

func sigToWitsig(sig *btcec.Signature) []byte {
	return append(sig.Serialize(), byte(SigHashType))
}

func decodeRawSignature(b []byte) *btcec.Signature {
	sig := &btcec.Signature{}
	sig.R = new(big.Int).SetBytes(b[:32])
	sig.S = new(big.Int).SetBytes(b[32:])
	return sig
}
