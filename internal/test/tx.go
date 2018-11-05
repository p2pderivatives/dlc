package test

import (
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

// TxVersion is a default tx version
const TxVersion = int32(2)

// SourceTx creates a transaction that has a dummy txin.
func SourceTx() *wire.MsgTx {
	tx := wire.NewMsgTx(TxVersion)
	prevOut := wire.NewOutPoint(&chainhash.Hash{}, ^uint32(0))
	txIn := wire.NewTxIn(prevOut, []byte{txscript.OP_0, txscript.OP_0}, nil)
	tx.AddTxIn(txIn)
	return tx
}
