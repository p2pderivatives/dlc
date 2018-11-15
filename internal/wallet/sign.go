package wallet

import (
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/dgarage/dlc/internal/script"
)

// WitnessSignature returns witness signature
// by creating a new keyset and signing tx with its privkey
func (w *wallet) WitnessSignature(
	tx *wire.MsgTx, idx int, amt btcutil.Amount, sc []byte, pub *btcec.PublicKey,
) ([]byte, error) {
	mpaddr, err := w.managedPubKeyAddressFromPubkey(pub)
	if err != nil {
		return nil, err
	}

	priv, err := mpaddr.PrivKey()
	if err != nil {
		return nil, err
	}
	sign, err := script.WitnessSignature(tx, idx, int64(amt), sc, priv)
	return sign, err
}

func (w *wallet) WitnessSignatureWithCallback(
	tx *wire.MsgTx, idx int, amt btcutil.Amount, sc []byte, pub *btcec.PublicKey,
	privkeyConverter PrivateKeyConverter,
) (sign []byte, err error) {
	return
}
