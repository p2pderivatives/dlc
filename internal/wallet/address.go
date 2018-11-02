package wallet

import (
	"fmt"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcwallet/waddrmgr"
	"github.com/btcsuite/btcwallet/walletdb"
	"github.com/dgarage/dlc/internal/script"
)

func (w *wallet) NewPubkey() (pub *btcec.PublicKey, err error) {

	testPrivPass := []byte("81lUHXnOMZ@?XXd7O9xyDIWIbXX-lj")

	mAddrs, err := w.newAddress(waddrmgr.KeyScopeBIP0084, testPrivPass, uint32(1), uint32(1))
	if err != nil {
		return nil, err
	}
	// pub = (waddrmgr.ManagedPubKeyAddress(mAddr[0])).PubKey()
	fmt.Printf("MADDRS[0}\n%+v\n", mAddrs[0])

	pub = (mAddrs[0].(waddrmgr.ManagedPubKeyAddress)).PubKey()
	fmt.Printf("PUB\n%+v\n", pub)

	return pub, err
}

func (w *wallet) NewWitnessPubkeyScript() (pkScript []byte, err error) {
	var pub *btcec.PublicKey
	pub, err = w.NewPubkey()
	if err != nil {
		return
	}
	return script.P2WPKHpkScript(pub)
}

// NewAddress returns a new ManagedAddress for a given scope and account number.
// NOTE: this function callsNextExternalAddresses to generate a ManagadAdddress.
func (w *wallet) newAddress(scope waddrmgr.KeyScope, privPass []byte,
	account uint32, numAddresses uint32) ([]waddrmgr.ManagedAddress, error) {
	// unlock Manager
	err := walletdb.Update(w.db, func(tx walletdb.ReadWriteTx) error {
		ns := tx.ReadWriteBucket(waddrmgrNamespaceKey)
		e := w.manager.Unlock(ns, privPass)
		return e
	})
	if err != nil {
		return nil, err
	}

	// get ScopedKeyManager
	scopedMgr, err := w.manager.FetchScopedKeyManager(scope)
	if err != nil {
		return nil, err
	}

	var addrs []waddrmgr.ManagedAddress
	err = walletdb.Update(w.db, func(tx walletdb.ReadWriteTx) error {
		ns := tx.ReadWriteBucket(waddrmgrNamespaceKey)
		var e error
		addrs, e = scopedMgr.NextExternalAddresses(ns, account, numAddresses)
		return e
	})
	if err != nil {
		return nil, err
	}

	return addrs, nil
}
