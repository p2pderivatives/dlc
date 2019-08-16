package dlc

import (
	"errors"
	"testing"

	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/p2pderivatives/dlc/internal/mocks/walletmock"
	"github.com/p2pderivatives/dlc/internal/test"
	"github.com/p2pderivatives/dlc/pkg/wallet"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// PrepareFundTx should fail if the party doesn't have enough balance
func TestPrepareFundTxNotEnoughUtxos(t *testing.T) {
	setupWallet := func() *walletmock.Wallet {
		w := setupTestWallet()
		w.On("SelectUnspent",
			mock.Anything, mock.Anything, mock.Anything,
		).Return(
			[]wallet.Utxo{}, btcutil.Amount(0), errors.New("not enough utxos"))
		return w
	}

	b := setupBuilder(FirstParty, setupWallet, newTestConditions)

	err := b.PrepareFundTx()
	assert.NotNil(t, err) // not enough balance for fee
}

// PrepareFundTx should prepare the txins and txouts of fundtx
func TestPrepareFundTx(t *testing.T) {
	assert := assert.New(t)

	setupWallet := func() *walletmock.Wallet {
		return mockSelectUnspent(
			setupTestWallet(), 1, 1, nil)
	}

	b := setupBuilder(
		FirstParty, setupWallet, newTestConditions)

	err := b.PrepareFundTx()
	assert.Nil(err)

	chaddr := b.Contract.ChangeAddrs[b.party]
	assert.NotEmpty(chaddr, "change address")
	utxos := b.Contract.Utxos[b.party]
	assert.NotNil(utxos, "utxos")

	// TODO: check if total amount is enough
}

// PrepareFundTx shouldn't have txouts if no changes
func TestPrepareFundTxNoChange(t *testing.T) {
	assert := assert.New(t)

	setupWallet := func() *walletmock.Wallet {
		return mockSelectUnspent(
			setupTestWallet(), 1, 0, nil)
	}

	b := setupBuilder(FirstParty, setupWallet, newTestConditions)

	err := b.PrepareFundTx()
	assert.Nil(err)

	utxos := b.Contract.Utxos[b.party]
	assert.NotNil(utxos, "utxos")
}

func testFundTx(t *testing.T, conditions func() *Conditions, expectedOutLen int) {
	assert := assert.New(t)

	// init builders
	b1 := setupBuilder(FirstParty, setupTestWallet, conditions)
	b2 := setupBuilder(SecondParty, setupTestWallet, conditions)

	// prep
	stepPrepare(b1)
	stepPrepare(b2)

	// fail if it hasn't received a pubkey from the counterparty
	_, err := b1.Contract.FundTx()
	assert.NotNil(err)

	// receive pubkey and utxos and change address from the counterparty
	stepSendRequirments(b2, b1)

	tx, err := b1.Contract.FundTx()
	assert.Nil(err)
	assert.Len(tx.TxIn, 2)  // funds from both parties
	assert.Len(tx.TxOut, expectedOutLen)
}

func TestFundTxNoPremium(t *testing.T) {
	expectedOutLen := 3  // 1 for redeemtx and 2 for changes
	testFundTx(t, newTestConditions, expectedOutLen)
}

func TestFundTxWithPremium(t *testing.T) {
	expectedOutLen := 4  // 1 for redeemtx, 2 for changes and one for premium
	testFundTx(t, newTestConditionsWithPremium, expectedOutLen)
}

func TestRedeemFundTx(t *testing.T) {
	assert := assert.New(t)

	// init builders
	b1 := setupBuilder(FirstParty, setupTestWallet, newTestConditions)
	b2 := setupBuilder(SecondParty, setupTestWallet, newTestConditions)

	// preparation
	stepPrepare(b1)
	stepPrepare(b2)
	stepSendRequirments(b2, b1)
	stepSendRequirments(b1, b2)

	d := b1.Contract

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
