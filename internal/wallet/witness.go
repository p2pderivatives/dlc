package wallet

import (
	"fmt"

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

// SignP2WPKH appends witness to P2WPKH txin
func (wallet *Wallet) SignP2WPKH(
	tx *wire.MsgTx,
	idx int,
	amt int64,
	script []byte,
	pub *btcec.PublicKey,
) error {
	priv, err := wallet.privkeyFromPubkey(pub)
	if err != nil {
		return err
	}

	sign, err := witnessSignature(tx, idx, amt, script, priv)
	if err != nil {
		return err
	}

	wt := wire.TxWitness{sign, pub.SerializeCompressed()}
	tx.TxIn[idx].Witness = wt

	return nil
}

// WitnessSignature returns signature for given script
func witnessSignature(
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

// TODO: implement me
func (wallet *Wallet) privkeyFromPubkey(pub *btcec.PublicKey) (*btcec.PrivateKey, error) {
	return &btcec.PrivateKey{}, fmt.Errorf("Impelment me")
}
