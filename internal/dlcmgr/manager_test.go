package dlcmgr

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcwallet/walletdb"
	_ "github.com/btcsuite/btcwallet/walletdb/bdb" // blank import for bolt db driver
	"github.com/p2pderivatives/dlc/internal/test"
	"github.com/p2pderivatives/dlc/pkg/dlc"
	"github.com/stretchr/testify/assert"
)

var (
	testDBPrefix = "dlcmgr_"
	testDBName   = "dlcmgr.db"
)

func TestCreateAndOpen(t *testing.T) {
	assert := assert.New(t)

	db, closeFunc := newWalletDB()
	defer closeFunc()

	manager, err := Create(db)
	assert.NoError(err)
	assert.NotNil(manager)

	manager, err = Open(db)
	assert.NoError(err)
	assert.NotNil(manager)
}

func TestStoreContract(t *testing.T) {
	assert := assert.New(t)

	// create new manager
	db, closeFunc := newWalletDB()
	defer closeFunc()
	manager, _ := Create(db)

	key := []byte("testdlc")
	dOrig := newDLC()

	err := manager.StoreContract(key, dOrig)
	if assert.NoError(err) {
		d, err := manager.RetrieveContract(key)
		assert.NoError(err)
		assert.NotNil(d)
		assert.Equal(dOrig, d)
	}
}

func TestRetrieveContractNotExists(t *testing.T) {
	assert := assert.New(t)

	// create new manager
	db, closeFunc := newWalletDB()
	defer closeFunc()
	manager, _ := Create(db)

	key := []byte("not_exists")
	d, err := manager.RetrieveContract(key)
	assert.Nil(d)
	assert.Error(err)
	assert.IsType(err, &ContractNotExistsError{})
}

func newWalletDB() (walletdb.DB, func()) {
	path := testDBPath()
	db, _ := walletdb.Create("bdb", path)
	closeFunc := func() {
		if db != nil {
			db.Close()
		}
		os.RemoveAll(path)
	}
	return db, closeFunc
}

func testDBPath() string {
	dir, _ := ioutil.TempDir("", testDBPrefix)
	path := filepath.Join(dir, testDBName)
	return path
}

func newDLC() *dlc.DLC {
	conds := testConditions()
	n := len(conds.Deals)
	return &dlc.DLC{
		Conds:       conds,
		Oracle:      testOracle(n),
		Pubs:        testPubkeys(),
		Addrs:       testAddrs(),
		ChangeAddrs: testAddrs(),
		Utxos:       testUtxos(),
		FundWits:    testFundWits(),
		RefundSigs:  testRefundSigs(),
		ExecSigs:    testExecSigs(),
	}
}

func testConditions() *dlc.Conditions {
	net := &chaincfg.RegressionNetParams
	ftime := testFixingTime()
	famt1, _ := btcutil.NewAmount(1)
	famt2, _ := btcutil.NewAmount(1)
	feerate := btcutil.Amount(10)
	refundlc := uint32(1)
	deals := newDeals()

	conds, _ := dlc.NewConditions(
		net, ftime, famt1, famt2, feerate, feerate, refundlc, deals, nil)

	return conds
}

func testFixingTime() time.Time {
	t := time.Now().AddDate(0, 0, 1)
	y, m, d := t.Date()
	return time.Date(y, m, d, 12, 0, 0, 0, time.UTC)
}

func newDeals() []*dlc.Deal {
	deals := []*dlc.Deal{}

	total := 5
	nMsg := 3
	for i := 0; i < total+1; i++ {
		amt1 := btcutil.Amount(i)
		amt2 := btcutil.Amount(total - i)
		msgs := [][]byte{}
		for j := 0; j < nMsg; j++ {
			msgs = append(msgs, []byte{byte(i)})
		}
		deal := dlc.NewDeal(amt1, amt2, msgs)
		deals = append(deals, deal)
	}

	return deals
}

func testOracle(n int) *dlc.Oracle {
	o := dlc.NewOracle(n)

	for i := 0; i < n; i++ {
		_, pub := test.RandKeys()
		o.Commitments[i] = pub
	}

	o.Sig = []byte{1}
	o.SignedMsgs = [][]byte{{1}}
	return o
}

func testPubkeys() map[dlc.Contractor]*btcec.PublicKey {
	pubs := make(map[dlc.Contractor]*btcec.PublicKey)
	_, pub1 := test.RandKeys()
	_, pub2 := test.RandKeys()
	pubs[dlc.FirstParty] = pub1
	pubs[dlc.SecondParty] = pub2
	return pubs
}

func testAddrs() map[dlc.Contractor]btcutil.Address {
	randAddr := func() btcutil.Address {
		_, pub := test.RandKeys()
		sc := btcutil.Hash160(pub.SerializeCompressed())
		net := &chaincfg.RegressionNetParams
		addr, _ := btcutil.NewAddressWitnessPubKeyHash(sc, net)
		return addr
	}

	addrs := make(map[dlc.Contractor]btcutil.Address)
	addrs[dlc.FirstParty] = randAddr()
	addrs[dlc.SecondParty] = randAddr()

	return addrs
}

func testUtxos() map[dlc.Contractor][]*dlc.Utxo {
	randUtxos := func() []*dlc.Utxo {
		return []*dlc.Utxo{
			{
				TxID:         "",
				Vout:         1,
				Address:      "",
				Account:      "",
				ScriptPubKey: "",
				RedeemScript: "",
				Amount:       1,
				Spendable:    true,
			}}
	}

	utxos := make(map[dlc.Contractor][]*dlc.Utxo)
	utxos[dlc.FirstParty] = randUtxos()
	utxos[dlc.SecondParty] = randUtxos()
	return utxos
}

func testFundWits() map[dlc.Contractor][]wire.TxWitness {
	wits := make(map[dlc.Contractor][]wire.TxWitness)
	wit1 := [][]byte{{1}}
	wits[dlc.FirstParty] = []wire.TxWitness{wit1}
	wit2 := [][]byte{{1}}
	wits[dlc.SecondParty] = []wire.TxWitness{wit2}
	return wits
}

func testRefundSigs() map[dlc.Contractor][]byte {
	sigs := make(map[dlc.Contractor][]byte)
	sigs[dlc.FirstParty] = []byte{1}
	sigs[dlc.SecondParty] = []byte{1}
	return sigs
}

func testExecSigs() [][]byte {
	sigs := [][]byte{}
	sigs = append(sigs, []byte{1})
	sigs = append(sigs, []byte{2})
	return sigs
}
