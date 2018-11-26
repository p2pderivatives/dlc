package script

import (
	"testing"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/wire"
	"github.com/dgarage/dlc/internal/test"
	"github.com/dgarage/dlc/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestCEScript(t *testing.T) {
	assert := assert.New(t)

	priva, puba := test.RandKeys()
	privb, pubb := test.RandKeys()
	privm, pubm := test.RandKeys()
	amt := int64(10000)

	script, err := ContractExecutionScript(puba, pubb, pubm)
	assert.Nil(err)
	pkScript, err := P2WSHpkScript(script)
	assert.Nil(err)

	// prepare source tx
	sourceTx := test.NewSourceTx()
	sourceTx.AddTxOut(wire.NewTxOut(amt, pkScript))

	// create redeem tx
	redeemTx := test.NewRedeemTx(sourceTx, 0)

	// unlock with message sign
	privam, _ := btcec.PrivKeyFromBytes(
		btcec.S256(),
		utils.AddBigInts(priva.D, privm.D).Bytes())
	signam, err := WitnessSignature(redeemTx, 0, amt, script, privam)
	assert.Nil(err)
	redeemTx.TxIn[0].Witness = WitnessForCEScript(signam, script)
	err = test.ExecuteScript(pkScript, redeemTx, amt)
	assert.Nil(err)

	// unlock after delay
	redeemTx.TxIn[0].Sequence = ContractExecutionDelay
	signb, err := WitnessSignature(redeemTx, 0, amt, script, privb)
	assert.Nil(err)
	redeemTx.TxIn[0].Witness = WitnessForCEScriptAfterDelay(signb, script)
	err = test.ExecuteScript(pkScript, redeemTx, amt)
	assert.Nil(err)
}
