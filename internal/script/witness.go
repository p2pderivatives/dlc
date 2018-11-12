package script

import (
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
)

// P2WPKHpkScript creates a withenss script for given pubkey.
//
// ScriptCode:
//  OP_0 + HASH160(<public key>)
func P2WPKHpkScript(pub *btcec.PublicKey) ([]byte, error) {
	builder := txscript.NewScriptBuilder()
	builder.AddOp(txscript.OP_0)
	builder.AddData(btcutil.Hash160(pub.SerializeCompressed()))
	return builder.Script()
}

// P2WSHpkScript creates a witness script for given script.
//
// ScriptCode:
//  OP_0 + SHA256(script)
func P2WSHpkScript(script []byte) ([]byte, error) {
	builder := txscript.NewScriptBuilder()
	builder.AddOp(txscript.OP_0)
	builder.AddData(chainhash.HashB(script))
	return builder.Script()
}

// WitnessSignature returns a witness signature for given script
//
// TODO: Note that txscript.RawTxInWitnessSignature converts a script from p2wkh to p2pkh implicitly.
// https://github.com/btcsuite/btcd/blob/master/txscript/script.go#L488
// It's better to convert it on ourside explicitly.
func WitnessSignature(
	tx *wire.MsgTx, idx int, amt int64, script []byte, priv *btcec.PrivateKey,
) ([]byte, error) {
	sighash := txscript.NewTxSigHashes(tx)
	return txscript.RawTxInWitnessSignature(
		tx, sighash, idx, amt, script, txscript.SigHashAll, priv)
}
