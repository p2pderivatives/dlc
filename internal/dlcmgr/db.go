package dlcmgr

import (
	"errors"
	"fmt"

	"github.com/btcsuite/btcwallet/walletdb"
)

var (
	nsTop         = []byte("dlcmgr")
	nsContracts   = []byte("contracts")
	nsOracle      = []byte("oracle")
	nsNetParam    = []byte("net")
	nsConditions  = []byte("conds")
	nsPubkeys     = []byte("pubkeys")
	nsAddrs       = []byte("addrs")
	nsChangeAddrs = []byte("chaddrs")
	nsUtxos       = []byte("utxos")
	nsFundWits    = []byte("fundwits")
	nsRefundSigs  = []byte("refundsigs")
	nsExecSigs    = []byte("execsigs")
)

func createManager(db walletdb.DB) error {
	err := walletdb.Update(db, func(tx walletdb.ReadWriteTx) error {
		_, _, e := createBucketsIfNotExist(tx)
		return e
	})

	return err
}

func createBucketsIfNotExist(tx walletdb.ReadWriteTx) (
	walletdb.ReadWriteBucket, walletdb.ReadWriteBucket, error) {
	var err error
	top := tx.ReadWriteBucket(nsTop)
	if top == nil {
		top, err = tx.CreateTopLevelBucket(nsTop)
		if err != nil {
			return nil, nil, err
		}
	}

	contracts, err := top.CreateBucketIfNotExists(nsContracts)
	return top, contracts, err
}

func openManager(db walletdb.DB) *Manager {
	mgr := &Manager{db: db}
	return mgr
}

type BucketNotExistsError struct {
	error
}

func (m *Manager) updateContractBucket(
	k []byte, f func(walletdb.ReadWriteBucket) error) error {
	updateFunc := func(tx walletdb.ReadWriteTx) (e error) {
		// TODO: workaround for panicking inside callback function
		defer func() {
			if r := recover(); r != nil {
				e = r.(error)
			}
		}()
		_, contracts, e := createBucketsIfNotExist(tx)
		if e != nil {
			return e
		}
		bucket, e := contracts.CreateBucketIfNotExists(k)
		if e != nil {
			return e
		}
		return f(bucket)
	}
	err := walletdb.Update(m.db, updateFunc)
	return err
}

// ContractNotExistsError gets raised when contract doesn't exist
type ContractNotExistsError struct {
	error
}

func newContractNotExistsError(
	k []byte) *ContractNotExistsError {
	msg := fmt.Sprintf("Contract not exists. key: %s", k)
	return &ContractNotExistsError{error: errors.New(msg)}
}

func (m *Manager) viewContractBucket(
	k []byte, f func(walletdb.ReadBucket) error) error {
	viewFunc := func(tx walletdb.ReadTx) error {
		top := tx.ReadBucket(nsTop)
		if top == nil {
			msg := fmt.Sprintf("bucket doesn't exist. bucket name: %s", top)
			return BucketNotExistsError{error: errors.New(msg)}
		}
		contracts := top.NestedReadBucket(nsContracts)
		if contracts == nil {
			msg := fmt.Sprintf("bucket doesn't exist. bucket name: %s", nsContracts)
			return BucketNotExistsError{error: errors.New(msg)}
		}
		bucket := contracts.NestedReadBucket(k)
		if bucket == nil {
			return newContractNotExistsError(k)
		}
		return f(bucket)
	}
	return walletdb.View(m.db, viewFunc)
}
