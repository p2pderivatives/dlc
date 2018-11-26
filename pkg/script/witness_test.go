package script

import (
	"testing"

	"github.com/btcsuite/btcd/wire"
	"github.com/dgarage/dlc/internal/test"
	"github.com/stretchr/testify/assert"
)

func TestP2WPKHpkScript(t *testing.T) {
	assert := assert.New(t)

	priv, pub := test.RandKeys()
	amt := int64(10000)

	// create P2WPKHpkScript
	pkScript, err := P2WPKHpkScript(pub)
	assert.Nil(err)

	// prepare source tx
	sourceTx := test.NewSourceTx()
	sourceTx.AddTxOut(wire.NewTxOut(amt, pkScript))

	// create redeem tx
	redeemTx := test.NewRedeemTx(sourceTx, 0)

	// witness signature
	sign, err := WitnessSignature(redeemTx, 0, amt, pkScript, priv)
	assert.Nil(err)

	// redeem script
	wt := wire.TxWitness{sign, pub.SerializeCompressed()}
	redeemTx.TxIn[0].Witness = wt

	// execute script
	err = test.ExecuteScript(pkScript, redeemTx, amt)
	assert.Nil(err)
}

func TestMultiSigScript2of2(t *testing.T) {
	assert := assert.New(t)

	priv1, pub1 := test.RandKeys()
	priv2, pub2 := test.RandKeys()
	amt := int64(10000)

	script, err := FundScript(pub1, pub2)
	assert.Nil(err)
	pkScript, err := P2WSHpkScript(script)
	assert.Nil(err)

	// prepare source tx
	sourceTx := test.NewSourceTx()
	sourceTx.AddTxOut(wire.NewTxOut(amt, pkScript))

	// create redeem tx
	redeemTx := test.NewRedeemTx(sourceTx, 0)

	// witness signatures
	sign1, err := WitnessSignature(redeemTx, 0, amt, script, priv1)
	assert.Nil(err)
	sign2, err := WitnessSignature(redeemTx, 0, amt, script, priv2)
	assert.Nil(err)

	// redeem script
	wt := wire.TxWitness{[]byte{}, sign1, sign2, script}
	redeemTx.TxIn[0].Witness = wt

	// execute script
	err = test.ExecuteScript(pkScript, redeemTx, amt)
	assert.Nil(err)
}
