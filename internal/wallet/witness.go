package wallet

import (
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
)

// P2WPKHpkScript creates P2WPKH pk script
//   OP_0 + HASH160(<public key>)
func P2WPKHpkScript(pub *btcec.PublicKey) ([]byte, error) {
	builder := txscript.NewScriptBuilder()
	builder.AddOp(txscript.OP_0)
	builder.AddData(btcutil.Hash160(pub.SerializeCompressed()))
	return builder.Script()
}

// WitnessSignature returns signature for given script
func WitnessSignature(
	tx *wire.MsgTx,
	idx int,
	amt int64,
	script []byte,
	priv *btcec.PrivateKey,
) ([]byte, error) {
	sighash := txscript.NewTxSigHashes(tx)
	return txscript.RawTxInWitnessSignature(
		tx, sighash, idx, amt, script, txscript.SigHashAll, priv)
}

// WitnessForP2WPKH returns witness for P2WPKH
func WitnessForP2WPKH(sign []byte, pub *btcec.PublicKey) [][]byte {
	tw := wire.TxWitness{}
	tw = append(tw, sign)
	tw = append(tw, pub.SerializeCompressed())
	return tw
}
