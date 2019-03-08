package dlc

import (
	"bytes"
	"encoding/hex"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
)

// utxoToTxIn converts utxo to txin
func utxoToTxIn(utxo *Utxo) (*wire.TxIn, error) {
	txid, err := chainhash.NewHashFromStr(utxo.TxID)
	if err != nil {
		return nil, err
	}
	op := wire.NewOutPoint(txid, utxo.Vout)
	txin := wire.NewTxIn(op, nil, nil)
	return txin, nil
}

func txToHex(tx *wire.MsgTx) (string, error) {
	// Serialize the transaction and convert to hex string.
	buf := bytes.NewBuffer(make([]byte, 0, tx.SerializeSize()))
	if err := tx.Serialize(buf); err != nil {
		return "", err
	}
	h := hex.EncodeToString(buf.Bytes())
	return h, nil
}

func hexToTx(txHex string) (tx *wire.MsgTx, err error) {
	txbin, err := hex.DecodeString(txHex)
	if err != nil {
		return nil, err
	}
	bufr := bytes.NewReader(txbin)
	err = tx.Deserialize(bufr)
	return
}
