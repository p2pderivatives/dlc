package wallet

import (
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcwallet/waddrmgr"
	"github.com/btcsuite/btcwallet/walletdb"
)

func (w *wallet) NewPubkey() (pub *btcec.PublicKey, err error) {
	mAddrs, err := w.newAddress(uint32(1))
	if err != nil {
		return nil, err
	}
	pub = (mAddrs[0].(waddrmgr.ManagedPubKeyAddress)).PubKey()
	return pub, err
}

func (w *wallet) NewWitnessPubkeyScript() (pkScript []byte, err error) {
	var pub *btcec.PublicKey
	pub, err = w.NewPubkey()
	if err != nil {
		return
	}
	return P2WPKHpkScript(pub)
}

// P2WPKHpkScript creates a withenss script for given pubkey
func P2WPKHpkScript(pub *btcec.PublicKey) ([]byte, error) {
	builder := txscript.NewScriptBuilder()
	builder.AddOp(txscript.OP_0)
	builder.AddData(btcutil.Hash160(pub.SerializeCompressed()))
	return builder.Script()
}

// NewAddress returns a new ManagedAddress
// NOTE: this function calls NextExternalAddresses to generate a ManagadAdddress.
func (w *wallet) newAddress(
	numAddresses uint32) ([]waddrmgr.ManagedAddress, error) {
	scopedMgr, err := w.manager.FetchScopedKeyManager(waddrmgrKeyScope)
	if err != nil {
		return nil, err
	}

	var addrs []waddrmgr.ManagedAddress
	err = walletdb.Update(w.db, func(tx walletdb.ReadWriteTx) error {
		ns := tx.ReadWriteBucket(waddrmgrNamespaceKey)
		var e error
		addrs, e = scopedMgr.NextExternalAddresses(ns, w.account, numAddresses)
		return e
	})
	if err != nil {
		return nil, err
	}

	return addrs, nil
}
