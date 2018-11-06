package dlc

import (
	"errors"
	"testing"

	"github.com/btcsuite/btcutil"
	"github.com/dgarage/dlc/internal/mocks"
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
	_, pub := test.RandKeys()
	w.On("NewPubkey").Return(pub, nil)
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
	// first party
	w1 := setupTestWallet()
	b1 := NewBuilder(FirstParty, mockSelectUnspent(w1, 1, 1, nil))
	b1.SetFundAmounts(1, 1)
	b1.PrepareFundTxIns()
	b1.PrepareFundPubkey()

	w2 := setupTestWallet()
	b2 := NewBuilder(SecondParty, mockSelectUnspent(w2, 1, 1, nil))
	b2.SetFundAmounts(1, 1)
	b2.PrepareFundTxIns()
	b2.PrepareFundPubkey()

	b1.CopyFundTxReqsFromCounterparty(b2.DLC())
	d := b1.DLC()

	tx, err := d.FundTx()

	assert := assert.New(t)
	assert.Nil(err)
	assert.Len(tx.TxIn, 2)  // funds from both parties
	assert.Len(tx.TxOut, 3) // 1 for reddemtx and 2 for changes
}
