package wallet

import (
	"errors"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcwallet/waddrmgr"
	"github.com/btcsuite/btcwallet/walletdb"
)

// NewPubkey returns a new btcec.PublicKey type public key
func (w *Wallet) NewPubkey() (pub *btcec.PublicKey, err error) {
	mAddr, err := w.newAddress()
	if err != nil {
		return nil, err
	}
	mpka := mAddr.(waddrmgr.ManagedPubKeyAddress)
	pub = mpka.PubKey()
	return pub, err
}

// NewAddress creates a new address managed by wallet
func (w *Wallet) NewAddress() (btcutil.Address, error) {
	maddr, err := w.newAddress()
	if err != nil {
		return nil, err
	}

	return maddr.Address(), nil
}

func (w *Wallet) ImportAddress(addr btcutil.Address) error {
	var maddr waddrmgr.ManagedAddress
	walletdb.View(w.db, func(tx walletdb.ReadTx) (e error) {
		ns := tx.ReadBucket(waddrmgrNamespaceKey)
		maddr, e = w.manager.Address(ns, addr)
		return e
	})
	if maddr != nil {
		return nil
	}

	// generate new addresses until reaching to the target address
	for i := 0; i < 100; i++ {
		newAddr, err := w.NewAddress()
		if err != nil {
			return err
		}

		if newAddr.EncodeAddress() == addr.EncodeAddress() {
			return nil
		}
	}

	return errors.New("Failed to import address")
}

// newAddress returns a new ManagedAddress
// NOTE: this function calls NextExternalAddresses to generate a ManagadAdddress.
func (w *Wallet) newAddress() (waddrmgr.ManagedAddress, error) {
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

	addr := addrs[0]

	// register address to bitcoind
	err = w.rpc.ImportAddressRescan(addr.Address().EncodeAddress(), "", false)
	if err != nil {
		return nil, err
	}

	return addr, nil
}

func (w *Wallet) managedPubKeyAddressFromPubkey(
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
