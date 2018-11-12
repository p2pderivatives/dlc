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
	var amt1, amt2 btcutil.Amount
	amt1, amt2 = 1, 1
	dID1 := setupDeal(b, amt1, amt2)

	// fail without oracle's message commitment
	_, err := b.dlc.ContractExecutionTx(b.party, dID1)
	assert.NotNil(err)

	// set message commitment
	_, msgCommit1 := test.RandKeys()
	b.SetMsgCommitmentToDeal(dID1, msgCommit1)

	// txout should have 2 entries
	tx1, err := b.dlc.ContractExecutionTx(b.party, dID1)
	assert.Nil(err)
	assert.Len(tx1.TxOut, 2)

	// A deal that destibutes to only one party
	amt1, amt2 = 2, 0
	dID2 := setupDeal(b, amt1, amt2)
	_, msgCommit2 := test.RandKeys()
	b.SetMsgCommitmentToDeal(dID2, msgCommit2)

	// txout should have only 1 entry
	tx2, err := b.dlc.ContractExecutionTx(b.party, dID2)
	assert.Nil(err)
	assert.Len(tx2.TxOut, 1)
}

func TestSignedContractExecutionTx(t *testing.T) {
	assert := assert.New(t)

	// setup
	b1, b2 := setupContractors()
	_, msgCommit := test.RandKeys()
	dID := setupDeal(b1, 1, 1)
	b1.SetMsgCommitmentToDeal(dID, msgCommit)
	_ = setupDeal(b2, 1, 1)
	b2.SetMsgCommitmentToDeal(dID, msgCommit)

	// fail without the counterparty's sign
	var err error
	_, err = b1.SignedContractExecutionTx(dID)
	assert.NotNil(err)
	_, err = b2.SignedContractExecutionTx(dID)
	assert.NotNil(err)

	// exchange signs
	sign1, err := b1.SignContractExecutionTx(dID)
	assert.Nil(err)
	sign2, err := b2.SignContractExecutionTx(dID)
	assert.Nil(err)

	err = b1.SetContextExecutionSign(dID, sign2)
	assert.Nil(err)
	err = b2.SetContextExecutionSign(dID, sign1)
	assert.Nil(err)

	// no errors with the counterparty's sign
	tx1, err := b1.SignedContractExecutionTx(dID)
	assert.Nil(err)
	tx2, err := b2.SignedContractExecutionTx(dID)
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
	// init first party
	w1 := setupTestWallet()
	b1 = NewBuilder(FirstParty, mockSelectUnspent(w1, 1, 1, nil))
	b1.SetFundAmounts(1, 1)
	b1.PreparePubkey()
	b1.PrepareFundTxIns()

	// init second party
	w2 := setupTestWallet()
	b2 = NewBuilder(SecondParty, mockSelectUnspent(w2, 1, 1, nil))
	b2.SetFundAmounts(1, 1)
	b2.PreparePubkey()
	b2.PrepareFundTxIns()

	// exchange pubkeys
	b1.CopyReqsFromCounterparty(b2.DLC())
	b2.CopyReqsFromCounterparty(b1.DLC())

	return b1, b2
}

func setupDeal(b *Builder, amt1, amt2 btcutil.Amount) int {
	msgs := [][]byte{{1}, {1}}
	deal := NewDeal(amt1, amt2, msgs)
	return b.AddDeal(deal)
}

func runFundScript(b *Builder, tx *wire.MsgTx) error {
	d := b.DLC()
	fundtx, _ := d.FundTx()
	fout := fundtx.TxOut[fundTxOutAt]
	return test.ExecuteScript(fout.PkScript, tx, fout.Value)
}
