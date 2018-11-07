package test

import (
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

// TxVersion is a default tx version
const TxVersion = int32(2)

// NewSourceTx creates a transaction that has a dummy txin.
func NewSourceTx() *wire.MsgTx {
	tx := wire.NewMsgTx(TxVersion)
	prevOut := wire.NewOutPoint(&chainhash.Hash{}, ^uint32(0))
	txIn := wire.NewTxIn(prevOut, []byte{txscript.OP_0, txscript.OP_0}, nil)
	tx.AddTxIn(txIn)
	return tx
}

// NewRedeemTx creates a redeem tx using a given source tx
func NewRedeemTx(sourceTx *wire.MsgTx, index uint32) *wire.MsgTx {
	txHash := sourceTx.TxHash()
	outPt := wire.NewOutPoint(&txHash, index)
	tx := wire.NewMsgTx(TxVersion)
	tx.AddTxIn(wire.NewTxIn(outPt, nil, nil))
	return tx
}

// ExecuteScript runs given script
func ExecuteScript(pkScript []byte, tx *wire.MsgTx, amt int64) error {
	flags := txscript.StandardVerifyFlags
	vm, err := txscript.NewEngine(pkScript, tx, 0, flags, nil, nil, amt)
	if err != nil {
		return err
	}
	return vm.Execute()
}
