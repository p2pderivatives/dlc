// Package wallet project wallet.go
package wallet

import (
	"log"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/hdkeychain"
)

// Wallet is wallet
type Wallet struct {
	extKey *hdkeychain.ExtendedKey
	params chaincfg.Params
	size   int
	// rpc    *rpc.BtcRPC
	infos []*Info
}

// Info is info data.
type Info struct {
	idx uint32
	pub *btcec.PublicKey
	adr string
}

// NewWallet returns a new Wallet
// func NewWallet(params chaincfg.Params, rpc *rpc.BtcRPC, seed []byte) (*Wallet, error) {
func NewWallet(params chaincfg.Params, seed []byte) (*Wallet, error) {
	wallet := &Wallet{}
	wallet.params = params
	// wallet.rpc = rpc
	wallet.size = 16
	mExtKey, err := hdkeychain.NewMaster(seed, &params)
	if err != nil {
		log.Printf("hdkeychain.NewMaster error : %v", err)
		return nil, err
	}
	key := mExtKey
	// m/44'/coin-type'/0'/0
	path := []uint32{44 | hdkeychain.HardenedKeyStart,
		params.HDCoinType | hdkeychain.HardenedKeyStart,
		0 | hdkeychain.HardenedKeyStart, 0}
	for _, i := range path {
		key, err = key.Child(i)
		if err != nil {
			log.Printf("key.Child error : %v", err)
			return nil, err
		}
	}
	wallet.extKey = key
	wallet.infos = []*Info{}
	for i := 0; i < wallet.size; i++ {
		key, _ := wallet.extKey.Child(uint32(i))
		pub, _ := key.ECPubKey()
		adr, _ := btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(pub.SerializeCompressed()), &wallet.params)
		info := &Info{uint32(i), pub, adr.EncodeAddress()}
		wallet.infos = append(wallet.infos, info)
		// _, err = rpc.Request("importaddress", adr.EncodeAddress(), "", false)
		if err != nil {
			return nil, err
		}
	}
	return wallet, nil
}
