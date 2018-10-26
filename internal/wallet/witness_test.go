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

	wallet := test.NewWallet()
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

	// sign P2WPKH
	err := wallet.SignP2WPKH(redeemTx, 0, amt, pkScript, pub)
	assert.Nil(err)

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
