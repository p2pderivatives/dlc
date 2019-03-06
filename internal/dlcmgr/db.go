package dlcmgr

import "github.com/btcsuite/btcwallet/walletdb"

var (
	nsTop        = []byte("dlcmgr")
	nsContracts  = []byte("contracts")
	nsConditions = []byte("conds")
)

func createManager(db walletdb.DB) error {
	err := walletdb.Update(db, func(tx walletdb.ReadWriteTx) error {
		ns, e := tx.CreateTopLevelBucket(nsTop)
		if e != nil {
			return e
		}

		if _, e = ns.CreateBucket(nsContracts); e != nil {
			return e
		}
		if _, e = ns.CreateBucket(nsConditions); e != nil {
			return e
		}
		return nil
	})

	return err
}

func openManager(db walletdb.DB) *Manager {
	mgr := &Manager{db: db}
	return mgr
}

func (m *Manager) updateBucket(
	bucketName []byte, f func(walletdb.ReadWriteBucket) error) error {
	updateFunc := func(tx walletdb.ReadWriteTx) error {
		ns := tx.ReadWriteBucket(nsTop)
		bucket := ns.NestedReadWriteBucket(bucketName)
		return f(bucket)
	}
	return walletdb.Update(m.db, updateFunc)
}

func (m *Manager) viewBucket(
	bucketName []byte, f func(walletdb.ReadBucket) error) error {
	viewFunc := func(tx walletdb.ReadTx) error {
		ns := tx.ReadBucket(nsTop)
		bucket := ns.NestedReadBucket(bucketName)
		return f(bucket)
	}
	return walletdb.View(m.db, viewFunc)
}
