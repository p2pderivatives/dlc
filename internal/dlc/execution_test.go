package dlc

import (
	"testing"

	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/dgarage/dlc/internal/test"
	"github.com/stretchr/testify/assert"
)

func TestCotractExecutionTx(t *testing.T) {
	assert := assert.New(t)

	b1, b2 := setupCountractors()

	// prepare a deal
	var amt1, amt2 btcutil.Amount = 1, 1
	msgs := [][]byte{[]byte{1}, []byte{1}}
	deal := NewDeal(amt1, amt2, msgs)

	dID := b1.AddDeal(deal)
	_ = b2.AddDeal(deal)

	// fail without oracle's message commitment
	_, err := b1.SignContractExecutionTx(SecondParty, dID)
	assert.NotNil(err)
	_, err = b2.SignContractExecutionTx(FirstParty, dID)
	assert.NotNil(err)

	// oracle's message commitment/sign
	_, msgCommit := test.RandKeys()

	// set message commitment
	err = b1.SetMsgCommitmentToDeal(dID, msgCommit)
	assert.Nil(err)
	err = b2.SetMsgCommitmentToDeal(dID, msgCommit)
	assert.Nil(err)

	// fail without the counterparty's sign
	_, err = b1.SignedContractExecutionTx(dID)
	assert.NotNil(err)
	_, err = b2.SignedContractExecutionTx(dID)
	assert.NotNil(err)

	// exchange signs
	sign1, err := b1.SignContractExecutionTx(SecondParty, dID)
	assert.Nil(err)
	sign2, err := b2.SignContractExecutionTx(FirstParty, dID)
	assert.Nil(err)

	err = b1.SetContextExecutionSign(dID, sign2)
	assert.Nil(err)
	err = b2.SetContextExecutionSign(dID, sign1)
	assert.Nil(err)

	tx1, err := b1.SignedContractExecutionTx(dID)
	assert.Nil(err)
	tx2, err := b2.SignedContractExecutionTx(dID)
	assert.Nil(err)
	assert.NotEqual(tx1, tx2)

	// run script
	err = runFundScript(b1, tx1)
	assert.Nil(err)
	err = runFundScript(b2, tx2)
	assert.Nil(err)
}

func setupCountractors() (b1, b2 *Builder) {
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

func runFundScript(b *Builder, tx *wire.MsgTx) error {
	d := b.DLC()
	fundtx, _ := d.FundTx()
	fout := fundtx.TxOut[fundTxOutAt]
	return test.ExecuteScript(fout.PkScript, tx, fout.Value)
}
