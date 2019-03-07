package dlcmgr

import (
	"encoding/json"

	"github.com/btcsuite/btcwallet/walletdb"
	"github.com/p2pderivatives/dlc/pkg/dlc"
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
	storeFunc := func(b walletdb.ReadWriteBucket) error {
		var e error
		if e = storeConditions(b, d.Conds); e != nil {
			return e
		}
		if e = storePublicKeys(b, d.PublicKeys()); e != nil {
			return e
		}
		return nil
	}
	return m.updateContractBucket(k, storeFunc)
}

func storeConditions(
	b walletdb.ReadWriteBucket, conds *dlc.Conditions) error {
	serializedConds, e := json.Marshal(conds)
	if e != nil {
		return e
	}
	return b.Put(nsConditions, serializedConds)
}

func storePublicKeys(
	b walletdb.ReadWriteBucket, pubs dlc.PublicKeys) error {
	serializedPubs, e := json.Marshal(pubs)
	if e != nil {
		return e
	}
	return b.Put(nsPubkeys, serializedPubs)
}

func (m *Manager) RetrieveContract(k []byte) (*dlc.DLC, error) {
	var d *dlc.DLC
	retrieveFunc := func(b walletdb.ReadBucket) error {
		conds, e := retrieveConditions(b)
		if e != nil {
			return e
		}

		d = dlc.NewDLC(conds)

		pubs, e := retrievePublicKeys(b)
		if e != nil {
			return e
		}
		if e = d.ParsePublicKeys(pubs); e != nil {
			return e
		}

		return e
	}
	err := m.viewContractBucket(k, retrieveFunc)
	return d, err
}

func retrieveConditions(b walletdb.ReadBucket) (*dlc.Conditions, error) {
	data := b.Get(nsConditions)
	conds := &dlc.Conditions{}
	e := json.Unmarshal(data, conds)
	return conds, e
}

func retrievePublicKeys(b walletdb.ReadBucket) (dlc.PublicKeys, error) {
	data := b.Get(nsPubkeys)
	pubs := make(dlc.PublicKeys)
	e := json.Unmarshal(data, &pubs)
	return pubs, e
}
