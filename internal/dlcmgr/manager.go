package dlcmgr

import (
	"encoding/json"
	"fmt"

	"github.com/btcsuite/btcwallet/walletdb"
	"github.com/dgarage/dlc/pkg/dlc"
)

type Manager struct {
	db walletdb.DB
}

func Create(db walletdb.DB) (*Manager, error) {
	err := createManager(db)
	if err != nil {
		return nil, err
	}
	return openManager(db), nil
}

func Open(db walletdb.DB) (*Manager, error) {
	return openManager(db), nil
}

func (m *Manager) Close() error {
	return m.db.Close()
}

func (m *Manager) StoreContract(k []byte, d *dlc.DLC) error {
	storeFunc := func(bucket walletdb.ReadWriteBucket) error {
		serializedConds, e := json.Marshal(d.Conds)
		if e != nil {
			return e
		}
		return bucket.Put(k, serializedConds)
	}
	return m.updateBucket(nsContracts, storeFunc)
}

func (m *Manager) RetrieveContract(k []byte) (*dlc.DLC, error) {
	var d *dlc.DLC
	retrieveFunc := func(bucket walletdb.ReadBucket) error {
		data := bucket.Get(k)

		conds := &dlc.Conditions{}
		e := json.Unmarshal(data, conds)
		fmt.Println(conds)

		d = &dlc.DLC{
			Conds: conds,
		}
		return e
	}
	err := m.viewBucket(nsContracts, retrieveFunc)
	return d, err
}

func (m *Manager) ListContracts() error {
	forEachFunc := func(k, v []byte) error {
		fmt.Println(k, v)
		return nil
	}
	listFunc := func(bucket walletdb.ReadBucket) error {
		return bucket.ForEach(forEachFunc)
	}
	return m.viewBucket(nsContracts, listFunc)
}
