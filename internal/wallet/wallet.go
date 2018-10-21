// Package wallet project wallet.go
package wallet

import (
	"fmt"
	"log"
	"math/big"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
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

	// TODO: change later, not safe for protection!!!
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

// GetWitnessSignature returns signature
func (w *Wallet) GetWitnessSignature(tx *wire.MsgTx, idx int, amt int64,
	script []byte, pub *btcec.PublicKey) ([]byte, error) {
	return w.GetWitnessSignaturePlus(tx, idx, amt, script, pub, nil)
}

// GetWitnessSignaturePlus returns signature for added private key
func (w *Wallet) GetWitnessSignaturePlus(tx *wire.MsgTx, idx int, amt int64,
	script []byte, pub *btcec.PublicKey, add *big.Int) ([]byte, error) {
	var pri *btcec.PrivateKey
	for _, info := range w.infos {
		if info.pub.IsEqual(pub) {
			key, _ := w.extKey.Child(info.idx)
			pri, _ = key.ECPrivKey()
		}
	}
	if pri == nil {
		return nil, fmt.Errorf("unknown public key %x", pub.SerializeCompressed())
	}
	if add != nil {
		num := new(big.Int).Mod(new(big.Int).Add(pri.D, add), btcec.S256().N)
		pri, _ = btcec.PrivKeyFromBytes(btcec.S256(), num.Bytes())
	}
	sighash := txscript.NewTxSigHashes(tx)
	sign, err := txscript.RawTxInWitnessSignature(tx, sighash, idx, amt, script, txscript.SigHashAll, pri)
	if err != nil {
		return nil, err
	}
	return sign, nil
}
