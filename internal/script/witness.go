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

// MultiSigScript2of2 is a 2-of-2 multisig script
//
// ScriptCode:
//  OP_2
//    <public key first party>
//    <public key second party>
//  OP_2
//  OP_CHECKMULTISIG
func MultiSigScript2of2(pub1, pub2 *btcec.PublicKey) (script []byte, err error) {
	builder := txscript.NewScriptBuilder()
	builder.AddOp(txscript.OP_2)
	builder.AddData(pub1.SerializeCompressed())
	builder.AddData(pub2.SerializeCompressed())
	builder.AddOp(txscript.OP_2)
	builder.AddOp(txscript.OP_CHECKMULTISIG)
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
