package dlc

import (
	"errors"
	"testing"

	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/p2pderivatives/dlc/internal/test"
	"github.com/p2pderivatives/dlc/pkg/wallet"
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
	builder := NewBuilder(FirstParty, testWallet, conds, regtestNetParams)

	err := builder.PrepareFundTx()
	assert.NotNil(t, err) // not enough balance for fee
}

// PrepareFundTx should prepare the txins and txouts of fundtx
func TestPrepareFundTx(t *testing.T) {
	assert := assert.New(t)

	conds := newTestConditions()

	// mock wallet
	testWallet := setupTestWallet()
	var balance, change btcutil.Amount = 1, 1
	mockSelectUnspent(testWallet, balance, change, nil)

	b := NewBuilder(FirstParty, testWallet, conds, regtestNetParams)

	err := b.PrepareFundTx()
	assert.Nil(err)

	chaddr := b.dlc.ChangeAddrs[b.party]
	assert.NotEmpty(chaddr, "change address")
	utxos := b.dlc.Utxos[b.party]
	assert.NotNil(utxos, "utxos")

	// TODO: check if total amount is enough
}

// PrepareFundTx shouldn't have txouts if no changes
func TestPrepareFundTxNoChange(t *testing.T) {
	assert := assert.New(t)

	// prepare mock wallet
	testWallet := setupTestWallet()
	var balance, change btcutil.Amount = 1, 0
	mockSelectUnspent(testWallet, balance, change, nil)

	conds := newTestConditions()
	b := NewBuilder(FirstParty, testWallet, conds, regtestNetParams)

	err := b.PrepareFundTx()
	assert.Nil(err)

	chaddr := b.dlc.ChangeAddrs[b.party]
	assert.Empty(chaddr, "change address")
	utxos := b.dlc.Utxos[b.party]
	assert.NotNil(utxos, "utxos")
}

func TestFundTx(t *testing.T) {
	assert := assert.New(t)
	conds := newTestConditions()

	// first party
	w1 := setupTestWallet()
	w1 = mockSelectUnspent(w1, 1000, 1, nil)
	b1 := NewBuilder(FirstParty, w1, conds, regtestNetParams)
	b1.PrepareFundTx()
	b1.PreparePubkey()

	// second party
	w2 := setupTestWallet()
	w2 = mockSelectUnspent(w2, 1000, 1, nil)
	b2 := NewBuilder(SecondParty, w2, conds, regtestNetParams)
	b2.PrepareFundTx()
	b2.PreparePubkey()

	// fail if it hasn't received a pubkey from the counterparty
	d := b1.DLC()
	_, err := d.FundTx()
	assert.NotNil(err)

	// receive pubkey and utxos and change address from the counterparty
	stepSendRequirments(b2, b1)

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
	w1 = mockSelectUnspent(w1, 1000, 1, nil)
	b1 := NewBuilder(FirstParty, w1, conds, regtestNetParams)
	b1.PreparePubkey()
	b1.PrepareFundTx()

	// init second party
	w2 := setupTestWallet()
	w2 = mockSelectUnspent(w2, 1000, 1, nil)
	b2 := NewBuilder(SecondParty, w2, conds, regtestNetParams)
	b2.PreparePubkey()
	b2.PrepareFundTx()

	// exchange pubkey, utxos, change address
	stepSendRequirments(b2, b1)
	stepSendRequirments(b1, b2)

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
