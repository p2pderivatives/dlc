package dlc

import (
	"testing"

	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/dgarage/dlc/internal/test"
	"github.com/stretchr/testify/assert"
)

func TestContractExecutionTx(t *testing.T) {
	assert := assert.New(t)

	b, _ := setupContractors()

	// A deal that has both amounts are > 0
	var amt1, amt2 btcutil.Amount = 1, 1
	dID, deal := setupDeal(b, amt1, amt2)

	// fail without oracle's message commitment
	_, err := b.dlc.ContractExecutionTx(b.party, deal)
	assert.NotNil(err)

	// set message commitment
	_, msgCommit := test.RandKeys()
	b.SetMsgCommitmentToDeal(dID, msgCommit)

	// txout should have 2 entries
	tx, err := b.dlc.ContractExecutionTx(b.party, deal)
	assert.Nil(err)
	assert.Len(tx.TxOut, 2)
	assert.Equal(int64(amt1), tx.TxOut[0].Value)
	assert.Equal(int64(amt2), tx.TxOut[1].Value)
}

// An edge case that a executing party tx takes all funds
func TestContractExecutionTxTakeAll(t *testing.T) {
	b, _ := setupContractors()

	var amt1, amt2 btcutil.Amount = 1, 0
	dID, deal := setupDeal(b, amt1, amt2)
	_, msgCommit := test.RandKeys()
	b.SetMsgCommitmentToDeal(dID, msgCommit)

	tx, err := b.dlc.ContractExecutionTx(b.party, deal)

	// asserions
	assert := assert.New(t)
	assert.Nil(err)
	assert.Len(tx.TxOut, 1)
	assert.Equal(int64(amt1), tx.TxOut[0].Value)
}

// An edge case that a executing party tx takes nothing
func TestContractExecutionTxTakeNothing(t *testing.T) {
	b, _ := setupContractors()

	var amt1, amt2 btcutil.Amount = 0, 1
	dID, deal := setupDeal(b, amt1, amt2)
	_, msgCommit := test.RandKeys()
	b.SetMsgCommitmentToDeal(dID, msgCommit)

	tx, err := b.dlc.ContractExecutionTx(b.party, deal)

	// asserions
	assert := assert.New(t)
	assert.Nil(tx)
	assert.NotNil(err)
	assert.IsType(&CETTakeNothingError{}, err)
}

func TestFixDeal(t *testing.T) {
	b, _ := setupContractors()

	dID, _ := setupDeal(b, 1, 1)
	msgPriv, msgCommit := test.RandKeys()
	b.SetMsgCommitmentToDeal(dID, msgCommit)

	err := b.FixDeal(dID, msgPriv.D.Bytes())
	assert.Nil(t, err)
}

func TestSignedContractExecutionTx(t *testing.T) {
	assert := assert.New(t)

	// setup
	b1, b2 := setupContractors()
	msgPriv, msgCommit := test.RandKeys()
	dID, deal := setupDeal(b1, 1, 1)
	b1.SetMsgCommitmentToDeal(dID, msgCommit)
	_, _ = setupDeal(b2, 1, 1)
	b2.SetMsgCommitmentToDeal(dID, msgCommit)

	msign := msgPriv.D.Bytes()
	b1.FixDeal(dID, msign)
	b2.FixDeal(dID, msign)

	// fail without the counterparty's sign
	var err error
	_, err = b1.SignedContractExecutionTx()
	assert.NotNil(err)
	_, err = b2.SignedContractExecutionTx()
	assert.NotNil(err)

	// exchange signs
	sign1, err := b1.SignContractExecutionTx(deal)
	assert.Nil(err)
	sign2, err := b2.SignContractExecutionTx(deal)
	assert.Nil(err)

	err = b1.AcceptCETxSign(dID, sign2)
	assert.Nil(err)
	err = b2.AcceptCETxSign(dID, sign1)
	assert.Nil(err)

	// no errors with the counterparty's sign
	tx1, err := b1.SignedContractExecutionTx()
	assert.Nil(err)
	tx2, err := b2.SignedContractExecutionTx()
	assert.Nil(err)

	// each party has a tx that has the same txin but has different txouts
	assert.Len(tx1.TxOut, 2)
	assert.Len(tx2.TxOut, 2)
	assert.Equal(
		tx1.TxIn[fundTxInAt].PreviousOutPoint,
		tx2.TxIn[fundTxInAt].PreviousOutPoint)
	assert.Equal(tx1.TxOut[0].Value, tx2.TxOut[1].Value)
	assert.Equal(tx1.TxOut[1].Value, tx2.TxOut[0].Value)

	// Both parties are able to send the CET
	err = runFundScript(b1, tx1)
	assert.Nil(err)
	err = runFundScript(b2, tx2)
	assert.Nil(err)
}

func setupContractors() (b1, b2 *Builder) {
	conds, _ := NewConditions(1, 1, 1, 1, 1)

	// init first party
	w1 := setupTestWallet()
	w1 = mockSelectUnspent(w1, 1, 1, nil)
	b1 = NewBuilder(FirstParty, w1, conds)
	b1.PreparePubkey()
	b1.PrepareFundTxIns()

	// init second party
	w2 := setupTestWallet()
	w2 = mockSelectUnspent(w2, 1, 1, nil)
	b2 = NewBuilder(SecondParty, w2, conds)
	b2.PreparePubkey()
	b2.PrepareFundTxIns()

	// exchange pubkeys
	b1.CopyReqsFromCounterparty(b2.DLC())
	b2.CopyReqsFromCounterparty(b1.DLC())

	return b1, b2
}

func setupDeal(b *Builder, amt1, amt2 btcutil.Amount) (int, *Deal) {
	msgs := [][]byte{{1}, {1}}
	deal := NewDeal(amt1, amt2, msgs)
	idx := b.AddDeal(deal)
	return idx, deal
}

func runFundScript(b *Builder, tx *wire.MsgTx) error {
	d := b.DLC()
	fundtx, _ := d.FundTx()
	fout := fundtx.TxOut[fundTxOutAt]
	return test.ExecuteScript(fout.PkScript, tx, fout.Value)
}
