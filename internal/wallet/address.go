package wallet

import (
	"errors"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcwallet/waddrmgr"
	"github.com/btcsuite/btcwallet/walletdb"
)

// NewPubKey returns a new btcec.PublicKey type public key
func (w *wallet) NewPubkey() (pub *btcec.PublicKey, err error) {
	mAddr, err := w.newAddress()
	if err != nil {
		return nil, err
	}
	pub = (mAddr.(waddrmgr.ManagedPubKeyAddress)).PubKey()
	return pub, err
}

// NewAddress creates a new address managed by wallet
func (w *wallet) NewAddress() (btcutil.Address, error) {
	maddr, err := w.newAddress()
	if err != nil {
		return nil, err
	}

	return maddr.Address(), nil
}

// newAddress returns a new ManagedAddress
// NOTE: this function calls NextExternalAddresses to generate a ManagadAdddress.
func (w *wallet) newAddress() (waddrmgr.ManagedAddress, error) {
	scopedMgr, err := w.manager.FetchScopedKeyManager(waddrmgrKeyScope)
	if err != nil {
		return nil, err
	}

	var numAddresses uint32 = 1
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

	// register address to bitcoind
	mpka := addrs[0].(waddrmgr.ManagedPubKeyAddress)
	err = w.rpc.ImportAddress(mpka.ExportPubKey())
	if err != nil {
		return nil, err
	}

	return addrs[0], nil
}

func (w *wallet) managedPubKeyAddressFromPubkey(
	pub *btcec.PublicKey,
) (rmpaddr waddrmgr.ManagedPubKeyAddress, err error) {
	err = walletdb.View(w.db, func(tx walletdb.ReadTx) error {
		ns := tx.ReadBucket(waddrmgrNamespaceKey)
		if ns == nil {
			return errors.New("missing address manager namespace")
		}
		e := w.manager.ForEachActiveAccountAddress(ns, w.account,
			func(maddr waddrmgr.ManagedAddress) error {
				mpaddr, ok := maddr.(waddrmgr.ManagedPubKeyAddress)
				if !ok {
					return nil
				}
				if !mpaddr.PubKey().IsEqual(pub) {
					return nil
				}
				rmpaddr = mpaddr
				return nil
			})
		return e
	})
	if rmpaddr == nil {
		msg := "No pubkey address is found associated with the given pubkey"
		return nil, errors.New(msg)
	}
	return rmpaddr, err
}
