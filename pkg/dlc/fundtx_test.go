package dlc

import (
	"errors"
	"testing"

	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/dgarage/dlc/internal/test"
	"github.com/dgarage/dlc/pkg/wallet"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// PrepareFundTx should fail if the party doesn't have enough balance
func TestPrepareFundTxNotEnoughUtxos(t *testing.T) {
	testWallet := setupTestWallet()
	testWallet.On("SelectUnspent",
		mock.Anything, mock.Anything, mock.Anything,
	).Return(
		[]wallet.Utxo{}, btcutil.Amount(0), errors.New("not enough utxos"))

	conds := newTestConditions()
	builder := NewBuilder(FirstParty, testWallet, conds)

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

	conds := newTestConditions()
	b := NewBuilder(FirstParty, testWallet, conds)

	err := b.PrepareFundTxIns()
	assert.Nil(err)

	txins := b.dlc.FundTxReqs.txIns[b.party]
	assert.NotEmpty(txins, "txins")
	txout := b.dlc.FundTxReqs.txOut[b.party]
	assert.NotNil(txout, "txout")
}

// PrepareFundTx shouldn't have txouts if no changes
func TestPrepareFundTxNoChange(t *testing.T) {
	assert := assert.New(t)

	// prepare mock wallet
	testWallet := setupTestWallet()
	var balance, change btcutil.Amount = 1, 0
	mockSelectUnspent(testWallet, balance, change, nil)

	conds := newTestConditions()
	b := NewBuilder(FirstParty, testWallet, conds)

	err := b.PrepareFundTxIns()
	assert.Nil(err)

	txins := b.dlc.FundTxReqs.txIns[b.party]
	assert.NotEmpty(txins, "txins")
	txout := b.dlc.FundTxReqs.txOut[b.party]
	assert.Nil(txout, "txout")
}

func TestFundTx(t *testing.T) {
	assert := assert.New(t)
	conds := newTestConditions()

	// first party
	w1 := setupTestWallet()
	w1 = mockSelectUnspent(w1, 1, 1, nil)
	b1 := NewBuilder(FirstParty, w1, conds)
	b1.PrepareFundTxIns()
	b1.PreparePubkey()

	// second party
	w2 := setupTestWallet()
	w2 = mockSelectUnspent(w2, 1, 1, nil)
	b2 := NewBuilder(SecondParty, w2, conds)
	b2.PrepareFundTxIns()
	b2.PreparePubkey()

	// fail if it hasn't received a pubkey from the counterparty
	d := b1.DLC()
	_, err := d.FundTx()
	assert.NotNil(err)

	// receive pubkey from the counterparty
	b1.CopyReqsFromCounterparty(b2.DLC())

	d = b1.DLC()
	tx, err := d.FundTx()
	assert.Nil(err)
	assert.Len(tx.TxIn, 2)  // funds from both parties
	assert.Len(tx.TxOut, 3) // 1 for reddemtx and 2 for changes
}

func TestRedeemFundTx(t *testing.T) {
	assert := assert.New(t)
	conds := newTestConditions()

	// init first party
	w1 := setupTestWallet()
	w1 = mockSelectUnspent(w1, 1, 1, nil)
	b1 := NewBuilder(FirstParty, w1, conds)
	b1.PreparePubkey()
	b1.PrepareFundTxIns()

	// init second party
	w2 := setupTestWallet()
	w2 = mockSelectUnspent(w2, 1, 1, nil)
	b2 := NewBuilder(SecondParty, w2, conds)
	b2.PreparePubkey()
	b2.PrepareFundTxIns()

	// exchange pubkeys
	b1.CopyReqsFromCounterparty(b2.DLC())
	b2.CopyReqsFromCounterparty(b1.DLC())
	d := b1.DLC()

	// prepare redeem tx for testing. this will be a settlement tx or refund tx
	redeemtx, err := d.newRedeemTx()
	assert.Nil(err)

	// both parties signs redeem tx
	sig1, err := b1.witsigForFundScript(redeemtx)
	assert.Nil(err)
	sig2, err := b2.witsigForFundScript(redeemtx)
	assert.Nil(err)

	// create witness
	fsc, _ := d.fundScript()
	wt := wire.TxWitness{[]byte{}, sig1, sig2, fsc}
	redeemtx.TxIn[0].Witness = wt

	// run script
	fundtx, _ := d.FundTx()
	fout := fundtx.TxOut[fundTxOutAt]
	err = test.ExecuteScript(fout.PkScript, redeemtx, fout.Value)
	assert.Nil(err)
}
