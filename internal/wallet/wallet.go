// Package wallet project wallet.go
package wallet

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/btcsuite/btcwallet/waddrmgr"
	"github.com/btcsuite/btcwallet/walletdb"
	"github.com/btcsuite/btcwallet/wtxmgr"
)

const (
	PUBPASSPHRASE []byte = pubpass
	PRIPASSPHRASE []byte = pripass

	// BIRTHDAY  Time.time =
)

// Namespace bucket keys.
var (
	waddrmgrNamespaceKey = []byte("waddrmgr")
	wtxmgrNamespaceKey   = []byte("wtxmgr")
)

// Wallet is hierarchical deterministic wallet
type Wallet struct {
	extKey *hdkeychain.ExtendedKey
	params chaincfg.Params
	size   int
	// rpc    *rpc.BtcRPC
	pubKeyInfos []*PublicKeyInfo
	Manager     *waddrmgr.Manager
}

// PublicKeyInfo is publickey data.
type PublicKeyInfo struct {
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

	dbPath := filepath.Join(os.TempDir(), "dev.db")
	db, err := walletdb.Create("bdb", dbPath)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	err = walletdb.Update(db,
		func(tx walletdb.ReadWriteTx) error {
			addrmgrNs, err := tx.CreateTopLevelBucket(waddrmgrNamespaceKey)
			if err != nil {
				return err
			}
			txmgrNs, err := tx.CreateTopLevelBucket(wtxmgrNamespaceKey)
			if err != nil {
				return err
			}

			err = waddrmgr.Create(
				addrmgrNs, seed, pubPass, privPass, params, nil,
				birthday,
			)
			if err != nil {
				return err
			}
			return wtxmgr.Create(txmgrNs)
		})

	addrmgrNs, err := tx.CreateTopLevelBucket(waddrmgrNamespaceKey)
	if err != nil {
		return nil, err
	}

	wallet.Manager = waddrmgr.Create(ns, seed, PUBPASSPHRASE, PRIPASSPHRASE, params, nil, time.Now())

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
	wallet.pubKeyInfos = []*PublicKeyInfo{}
	for i := 0; i < wallet.size; i++ {
		key, _ := wallet.extKey.Child(uint32(i))
		pub, _ := key.ECPubKey()
		adr, _ := btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(pub.SerializeCompressed()), &wallet.params)
		info := &PublicKeyInfo{uint32(i), pub, adr.EncodeAddress()}
		wallet.pubKeyInfos = append(wallet.pubKeyInfos, info)
		// _, err = rpc.Request("importaddress", adr.EncodeAddress(), "", false)
		if err != nil {
			return nil, err
		}
	}
	return wallet, nil
}

func (w *Wallet) SendTx(tx *wire.MsgTx) (*chainhash.Hash, error) {
	allowHighFees := false
	//return w.rpc.SendRawTransaction(tx, allowHighFees)

	// testing
	// marshalled: `{"jsonrpc":"1.0","method":"sendrawtransaction","params":["1122"],"id":1}`,
	// unmarshalled: &btcjson.SendRawTransactionCmd{
	// 	HexTx:         "1122",
	// 	AllowHighFees: btcjson.Bool(false),
	// },
	// https://github.com/btcsuite/btcd/blob/fdfc19097e7ac6b57035062056f5b7b4638b8898/btcjson/chainsvrcmds_test.go#L903

}
