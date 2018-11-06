package wallet

import (
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/wire"
	"github.com/dgarage/dlc/internal/script"
)

// WitnessSignature returns witness signature
// by getting privkey from a given pubkey
func (w *wallet) WitnessSignature(
	tx *wire.MsgTx, idx int, amt int64, sc []byte, pub *btcec.PublicKey,
) (sign []byte, err error) {
	priv, err := w.privkeyFromPubkey(pub)
	if err != nil {
		return
	}

	return script.WitnessSignature(tx, idx, amt, sc, priv)
}

// privkeyFromPubkey retrieves a privkey for a given pubkey
func (w *wallet) privkeyFromPubkey(
	pub *btcec.PublicKey) (priv *btcec.PrivateKey, err error) {
	// TODO implement this function
	return
}
