package dlcmgr

import (
	"encoding/json"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcwallet/walletdb"
	"github.com/p2pderivatives/dlc/pkg/dlc"
)

// Manager manages contracts
type Manager struct {
	db walletdb.DB
}

// Create creates manager
func Create(db walletdb.DB) (*Manager, error) {
	err := createManager(db)
	if err != nil {
		return nil, err
	}
	return openManager(db), nil
}

// Open opens manager
func Open(db walletdb.DB) (*Manager, error) {
	return openManager(db), nil
}

// Close closes manager
func (m *Manager) Close() error {
	return m.db.Close()
}

// StoreContract persists DLC
func (m *Manager) StoreContract(k []byte, d *dlc.DLC) error {
	storeFunc := func(b walletdb.ReadWriteBucket) error {
		var e error
		if e = storeConditions(b, d.Conds); e != nil {
			return e
		}
		if e = storePublicKeys(b, d.PublicKeys()); e != nil {
			return e
		}
		if e = storeAddrs(b, d.Addresses()); e != nil {
			return e
		}
		if e = storeChangeAddrs(b, d.ChangeAddresses()); e != nil {
			return e
		}
		if e = storeUtxos(b, d.Utxos); e != nil {
			return e
		}
		if e = storeFundWits(b, d.FundWits); e != nil {
			return e
		}
		if e = storeRefundSigs(b, d.RefundSigs); e != nil {
			return e
		}
		if e = storeExecSigs(b, d.ExecSigs); e != nil {
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

func storeAddrs(
	b walletdb.ReadWriteBucket, addrs dlc.Addresses) error {
	serializedAddrs, e := json.Marshal(addrs)
	if e != nil {
		return e
	}
	return b.Put(nsAddrs, serializedAddrs)
}

func storeChangeAddrs(
	b walletdb.ReadWriteBucket, addrs dlc.Addresses) error {
	serializedAddrs, e := json.Marshal(addrs)
	if e != nil {
		return e
	}
	return b.Put(nsChangeAddrs, serializedAddrs)
}

func storeUtxos(
	b walletdb.ReadWriteBucket, utxos map[dlc.Contractor][]*dlc.Utxo) error {
	serializedUtxos, e := json.Marshal(utxos)
	if e != nil {
		return e
	}
	return b.Put(nsUtxos, serializedUtxos)
}

func storeFundWits(
	b walletdb.ReadWriteBucket, wits map[dlc.Contractor][]wire.TxWitness) error {
	serializedWits, e := json.Marshal(wits)
	if e != nil {
		return e
	}
	return b.Put(nsFundWits, serializedWits)
}

func storeRefundSigs(
	b walletdb.ReadWriteBucket, sigs map[dlc.Contractor][]byte) error {
	serializedSigs, e := json.Marshal(sigs)
	if e != nil {
		return e
	}
	return b.Put(nsRefundSigs, serializedSigs)
}

func storeExecSigs(
	b walletdb.ReadWriteBucket, sigs [][]byte) error {
	serializedSigs, e := json.Marshal(sigs)
	if e != nil {
		return e
	}
	return b.Put(nsExecSigs, serializedSigs)
}

// RetrieveContract retrieves stored DLC
func (m *Manager) RetrieveContract(k []byte) (*dlc.DLC, error) {
	var d *dlc.DLC
	retrieveFunc := func(b walletdb.ReadBucket) error {
		conds, e := retrieveConditions(b)
		if e != nil {
			return e
		}

		// TODO: store and retrieve netparam
		net := &chaincfg.RegressionNetParams
		d = dlc.NewDLC(conds, net)

		pubs, e := retrievePublicKeys(b)
		if e != nil {
			return e
		}
		if e = d.ParsePublicKeys(pubs); e != nil {
			return e
		}

		addrs, e := retrieveAddrs(b)
		if e != nil {
			return e
		}
		if e = d.ParseAddresses(addrs); e != nil {
			return e
		}

		chaddrs, e := retrieveChangeAddrs(b)
		if e != nil {
			return e
		}
		if e = d.ParseChangeAddresses(chaddrs); e != nil {
			return e
		}

		utxos, e := retrieveUtxos(b)
		if e != nil {
			return e
		}
		d.Utxos = utxos

		fundWits, e := retrieveFundWits(b)
		if e != nil {
			return e
		}
		d.FundWits = fundWits

		refundSigs, e := retrieveRefundSigs(b)
		if e != nil {
			return e
		}
		d.RefundSigs = refundSigs

		execSigs, e := retrieveExecSigs(b)
		if e != nil {
			return e
		}
		d.ExecSigs = execSigs

		return e
	}
	err := m.viewContractBucket(k, retrieveFunc)
	return d, err
}

func retrieveConditions(b walletdb.ReadBucket) (*dlc.Conditions, error) {
	data := b.Get(nsConditions)
	if len(data) == 0 {
		return nil, nil
	}
	conds := &dlc.Conditions{}
	e := json.Unmarshal(data, conds)
	return conds, e
}

func retrievePublicKeys(b walletdb.ReadBucket) (dlc.PublicKeys, error) {
	data := b.Get(nsPubkeys)
	if len(data) == 0 {
		return nil, nil
	}
	pubs := make(dlc.PublicKeys)
	e := json.Unmarshal(data, &pubs)
	if e != nil {
		return nil, e
	}
	return pubs, e
}

func retrieveAddrs(b walletdb.ReadBucket) (dlc.Addresses, error) {
	data := b.Get(nsAddrs)
	if len(data) == 0 {
		return nil, nil
	}
	addrs := make(dlc.Addresses)
	e := json.Unmarshal(data, &addrs)
	return addrs, e
}

func retrieveChangeAddrs(b walletdb.ReadBucket) (dlc.Addresses, error) {
	data := b.Get(nsChangeAddrs)
	if len(data) == 0 {
		return nil, nil
	}
	addrs := make(dlc.Addresses)
	e := json.Unmarshal(data, &addrs)
	return addrs, e
}

func retrieveUtxos(b walletdb.ReadBucket) (map[dlc.Contractor][]*dlc.Utxo, error) {
	data := b.Get(nsUtxos)
	if len(data) == 0 {
		return nil, nil
	}
	utxos := make(map[dlc.Contractor][]*dlc.Utxo)
	e := json.Unmarshal(data, &utxos)
	return utxos, e
}

func retrieveFundWits(b walletdb.ReadBucket) (map[dlc.Contractor][]wire.TxWitness, error) {
	data := b.Get(nsFundWits)
	if len(data) == 0 {
		return nil, nil
	}
	wits := make(map[dlc.Contractor][]wire.TxWitness)
	e := json.Unmarshal(data, &wits)
	return wits, e
}

func retrieveRefundSigs(b walletdb.ReadBucket) (map[dlc.Contractor][]byte, error) {
	data := b.Get(nsRefundSigs)
	if len(data) == 0 {
		return nil, nil
	}
	sigs := make(map[dlc.Contractor][]byte)
	e := json.Unmarshal(data, &sigs)
	return sigs, e
}

func retrieveExecSigs(b walletdb.ReadBucket) ([][]byte, error) {
	data := b.Get(nsExecSigs)
	if len(data) == 0 {
		return nil, nil
	}
	sigs := [][]byte{}
	e := json.Unmarshal(data, &sigs)
	return sigs, e
}
