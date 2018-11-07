package dlc

import (
	"errors"
	"testing"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/dgarage/dlc/internal/mocks"
	"github.com/dgarage/dlc/internal/script"
	"github.com/dgarage/dlc/internal/test"
	"github.com/dgarage/dlc/internal/wallet"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Hash of block 234439
var testTxID = "14a0810ac680a3eb3f82edc878cea25ec41d6b790744e5daeef"

// setup mocke wallet
func setupTestWallet() *mocks.Wallet {
	w := &mocks.Wallet{}
	priv, pub := test.RandKeys()
	w.On("NewPubkey").Return(pub, nil)
	w = mockWitnessSignature(w, pub, priv)
	return w
}

func mockSelectUnspent(
	w *mocks.Wallet, balance, change btcutil.Amount, err error) *mocks.Wallet {
	utxo := wallet.Utxo{
		TxID:   testTxID,
		Amount: float64(balance) / btcutil.SatoshiPerBitcoin,
	}
	w.On("SelectUnspent",
		mock.Anything, mock.Anything, mock.Anything,
	).Return([]wallet.Utxo{utxo}, change, err)

	return w
}

func mockWitnessSignature(
	w *mocks.Wallet, pub *btcec.PublicKey, priv *btcec.PrivateKey) *mocks.Wallet {
	call := w.On("WitnessSignature",
		mock.AnythingOfType("*wire.MsgTx"),
		mock.AnythingOfType("int"),
		mock.AnythingOfType("btcutil.Amount"),
		mock.AnythingOfType("[]uint8"),
		pub,
	)

	call.Run(func(args mock.Arguments) {
		tx := args.Get(0).(*wire.MsgTx)
		idx := args.Get(1).(int)
		amt := args.Get(2).(btcutil.Amount)
		sc := args.Get(3).([]uint8)
		sign, err := script.WitnessSignature(tx, idx, int64(amt), sc, priv)
		rargs := make([]interface{}, 2)
		rargs[0] = sign
		rargs[1] = err
		call.ReturnArguments = rargs
	})

	return w
}

// PrepareFundTx should fail if fund amounts aren't set
func TestPrepareFundTxNoFundAmounts(t *testing.T) {
	builder := NewBuilder(FirstParty, nil)

	err := builder.PrepareFundTxIns()
	assert.NotNil(t, err)
}

// PrepareFundTx should fail if the party doesn't have enough balance
func TestPrepareFundTxNotEnoughUtxos(t *testing.T) {
	testWallet := setupTestWallet()
	testWallet.On("SelectUnspent",
		mock.Anything, mock.Anything, mock.Anything,
	).Return(
		[]wallet.Utxo{}, btcutil.Amount(0), errors.New("not enough utxos"))

	builder := NewBuilder(FirstParty, testWallet)
	var famt btcutil.Amount = 1 // 1 satoshi
	builder.SetFundAmounts(famt, famt)

	err := builder.PrepareFundTxIns()
	assert.NotNil(t, err) // not enough balance for fee
}

// PrepareFundTx should prepare the txins and txouts of fundtx
func TestPrepareFundTx(t *testing.T) {
	assert := assert.New(t)

	// prepare mock wallet
	testWallet := setupTestWallet()
	var balance, change btcutil.Amount = 1, 1
	mockSelectUnspent(testWallet, balance, change, nil)

	b := NewBuilder(FirstParty, testWallet)
	b.SetFundAmounts(1, 1)

	err := b.PrepareFundTxIns()
	assert.Nil(err)

	txins := b.dlc.fundTxReqs.txIns[b.party]
	assert.NotEmpty(txins, "txins")
	txout := b.dlc.fundTxReqs.txOut[b.party]
	assert.NotNil(txout, "txout")
}

// PrepareFundTx shouldn't have txouts if no changes
func TestPrepareFundTxNoChange(t *testing.T) {
	assert := assert.New(t)

	// prepare mock wallet
	testWallet := setupTestWallet()
	var balance, change btcutil.Amount = 1, 0
	mockSelectUnspent(testWallet, balance, change, nil)

	b := NewBuilder(FirstParty, testWallet)
	b.SetFundAmounts(1, 1)

	err := b.PrepareFundTxIns()
	assert.Nil(err)

	txins := b.dlc.fundTxReqs.txIns[b.party]
	assert.NotEmpty(txins, "txins")
	txout := b.dlc.fundTxReqs.txOut[b.party]
	assert.Nil(txout, "txout")
}

func TestFundTx(t *testing.T) {
	assert := assert.New(t)

	// first party
	w1 := setupTestWallet()
	b1 := NewBuilder(FirstParty, mockSelectUnspent(w1, 1, 1, nil))
	b1.SetFundAmounts(1, 1)
	b1.PrepareFundTxIns()
	b1.PrepareFundPubkey()

	// second party
	w2 := setupTestWallet()
	b2 := NewBuilder(SecondParty, mockSelectUnspent(w2, 1, 1, nil))
	b2.SetFundAmounts(1, 1)
	b2.PrepareFundTxIns()
	b2.PrepareFundPubkey()

	// fail if it hasn't received a pubkey from the counterparty
	d := b1.DLC()
	_, err := d.FundTx()
	assert.NotNil(err)

	// receive pubkey from the counterparty
	b1.CopyFundTxReqsFromCounterparty(b2.DLC())

	d = b1.DLC()
	tx, err := d.FundTx()
	assert.Nil(err)
	assert.Len(tx.TxIn, 2)  // funds from both parties
	assert.Len(tx.TxOut, 3) // 1 for reddemtx and 2 for changes
}

func TestRedeemFundTx(t *testing.T) {
	assert := assert.New(t)

	// init first party
	w1 := setupTestWallet()
	b1 := NewBuilder(FirstParty, mockSelectUnspent(w1, 1, 1, nil))
	b1.SetFundAmounts(1, 1)
	b1.PrepareFundTxIns()
	b1.PrepareFundPubkey()

	// init second party
	w2 := setupTestWallet()
	b2 := NewBuilder(SecondParty, mockSelectUnspent(w2, 1, 1, nil))
	b2.SetFundAmounts(1, 1)
	b2.PrepareFundTxIns()
	b2.PrepareFundPubkey()

	// exchange pubkeys
	b1.CopyFundTxReqsFromCounterparty(b2.DLC())
	b2.CopyFundTxReqsFromCounterparty(b1.DLC())
	d := b1.DLC()

	// prepare redeem tx for testing. this will be a settlement tx or refund tx
	redeemtx, err := d.newRedeemTx()
	assert.Nil(err)

	// both parties signs redeem tx
	sign1, err := b1.witsigForRedeemTx(redeemtx)
	assert.Nil(err)
	sign2, err := b2.witsigForRedeemTx(redeemtx)
	assert.Nil(err)

	// create witness
	fsc, _ := d.fundScript()
	wt := wire.TxWitness{[]byte{}, sign1, sign2, fsc}
	redeemtx.TxIn[0].Witness = wt

	// run script
	pkScript, _ := script.P2WSHpkScript(fsc)
	famt, _ := d.fundAmount()
	err = test.ExecuteScript(pkScript, redeemtx, int64(famt))
	assert.Nil(err)
}
