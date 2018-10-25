package wallet

import (
	"testing"

	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/dgarage/dlc/internal/test"
	"github.com/stretchr/testify/assert"
)

func TestP2WPKHpkScript(t *testing.T) {
	var err error
	assert := assert.New(t)

	pri, pub := test.RandKeys()
	amt := int64(10000)
	pkScript, _ := P2WPKHpkScript(pub)

	// prepare source transaction
	sourceTx := test.SourceTx()

	// append P2WPKH script
	witOut := wire.NewTxOut(amt, pkScript)
	sourceTx.AddTxOut(witOut)

	// create redeem tx from source tx
	redeemTx := createRedeemTx(sourceTx)

	// append witness signature to redeem tx
	sign, err := WitnessSignature(redeemTx, 0, amt, pkScript, pri)
	assert.Nil(err)
	redeemTx.TxIn[0].Witness = WitnessForP2WPKH(sign, pub)

	// execute script
	err = executeScript(pkScript, redeemTx, amt)
	assert.Nil(err)
}

func createRedeemTx(sourceTx *wire.MsgTx) *wire.MsgTx {
	txHash := sourceTx.TxHash()
	outPt := wire.NewOutPoint(&txHash, 0)

	tx := wire.NewMsgTx(test.TxVersion)
	tx.AddTxIn(wire.NewTxIn(outPt, nil, nil))

	return tx
}

func executeScript(pkScript []byte, tx *wire.MsgTx, amt int64) error {
	flags := txscript.StandardVerifyFlags
	vm, err := txscript.NewEngine(pkScript, tx, 0, flags, nil, nil, amt)
	if err != nil {
		return err
	}

	return vm.Execute()
}
