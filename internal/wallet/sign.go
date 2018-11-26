package wallet

import (
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcwallet/waddrmgr"
	"github.com/btcsuite/btcwallet/walletdb"
	"github.com/dgarage/dlc/pkg/script"
	"github.com/dgarage/dlc/pkg/wallet"
)

// WitnessSignature returns witness signature
// by signing tx with the privkey of given pubkey
func (w *Wallet) WitnessSignature(
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

// WitnessSignatureWithCallback does the same with WitnessSignature does,
// but applying a func PrivateKeyConverter to raw privkey and use the result
// as privkey for calculating witness signature
func (w *Wallet) WitnessSignatureWithCallback(
	tx *wire.MsgTx, idx int, amt btcutil.Amount, sc []byte, pub *btcec.PublicKey,
	privkeyConverter wallet.PrivateKeyConverter,
) ([]byte, error) {
	mpaddr, err := w.managedPubKeyAddressFromPubkey(pub)
	if err != nil {
		return nil, err
	}

	rawpriv, err := mpaddr.PrivKey()
	if err != nil {
		return nil, err
	}
	priv, err := privkeyConverter(rawpriv)
	if err != nil {
		return nil, err
	}

	sign, err := script.WitnessSignature(tx, idx, int64(amt), sc, priv)
	return sign, err
}

// WitnessSignTxByIdxs returns witnesses associated to txins at given indices
func (w *Wallet) WitnessSignTxByIdxs(tx *wire.MsgTx, idxs []int) ([]wire.TxWitness, error) {
	wits := []wire.TxWitness{}
	for _, idx := range idxs {
		txin := tx.TxIn[idx]

		// txin -> utxo
		utxo, err := w.UtxoByTxIn(txin)
		if err != nil {
			return nil, err
		}

		// utxo -> managed address
		maddr, err := w.managedAddressByUtxo(utxo)
		if err != nil {
			return nil, err
		}

		// retrieve pubkey and privkey
		mpka := maddr.(waddrmgr.ManagedPubKeyAddress)
		pub := mpka.PubKey()
		priv, err := mpka.PrivKey()
		if err != nil {
			return nil, err
		}

		// calc witness signature
		amt, err := btcutil.NewAmount(utxo.Amount)
		if err != nil {
			return nil, err
		}

		sc, err := script.P2WPKHpkScript(pub)
		if err != nil {
			return nil, err
		}

		sign, err := script.WitnessSignature(tx, idx, int64(amt), sc, priv)
		if err != nil {
			return nil, err
		}

		// compose witness
		wit := wire.TxWitness{sign, pub.SerializeCompressed()}
		wits = append(wits, wit)
	}

	return wits, nil
}

// managedAddressByUtxo finds managed address by utxo
func (w *Wallet) managedAddressByUtxo(utxo wallet.Utxo) (maddr waddrmgr.ManagedAddress, err error) {
	onEachAddr := func(_maddr waddrmgr.ManagedAddress) error {
		if _maddr.Address().String() == utxo.Address {
			maddr = _maddr
		}
		return nil
	}
	onView := func(tx walletdb.ReadTx) error {
		return w.manager.ForEachActiveAccountAddress(
			tx.ReadBucket(waddrmgrNamespaceKey), w.account, onEachAddr)
	}
	err = walletdb.View(w.db, onView)
	if err != nil {
		return nil, err
	}

	if maddr == nil {
		errmsg := fmt.Sprintf(
			"managed address not found by utxo. utxo: %#v", utxo)
		return nil, errors.New(errmsg)
	}

	return maddr, nil
}
