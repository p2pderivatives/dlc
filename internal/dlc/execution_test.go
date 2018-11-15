package dlc

import (
	"testing"

	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/dgarage/dlc/internal/oracle"
	"github.com/dgarage/dlc/internal/test"
	"github.com/stretchr/testify/assert"
)

func TestContractExecutionTx(t *testing.T) {
	assert := assert.New(t)

	// A deal that has both amounts are > 0
	var damt1, damt2 btcutil.Amount = 1, 1
	b, _, dID, deal := setupContractorsUntilPubkeyExchange(damt1, damt2)

	// fail without oracle's message commitment
	_, err := b.dlc.ContractExecutionTx(b.party, deal, dID)
	assert.NotNil(err)

	// set message commitment
	_, C := test.RandKeys()
	b.dlc.oracleReqs.commitments[dID] = C

	// txout should have 2 entries
	tx, err := b.dlc.ContractExecutionTx(b.party, deal, dID)
	assert.Nil(err)
	assert.Len(tx.TxOut, 2)
	assert.Equal(int64(damt1), tx.TxOut[0].Value)
	assert.Equal(int64(damt2), tx.TxOut[1].Value)
}

// An edge case that a executing party tx takes all funds
func TestContractExecutionTxTakeAll(t *testing.T) {
	var damt1, damt2 btcutil.Amount = 1, 0
	b, _, dID, deal := setupContractorsUntilPubkeyExchange(damt1, damt2)
	_, C := test.RandKeys()
	b.dlc.oracleReqs.commitments[dID] = C

	tx, err := b.dlc.ContractExecutionTx(b.party, deal, dID)

	// asserions
	assert := assert.New(t)
	assert.Nil(err)
	assert.Len(tx.TxOut, 1)
	assert.Equal(int64(damt1), tx.TxOut[0].Value)
}

// An edge case that a executing party tx takes nothing
func TestContractExecutionTxTakeNothing(t *testing.T) {
	var damt1, damt2 btcutil.Amount = 0, 1
	b, _, dID, deal := setupContractorsUntilPubkeyExchange(damt1, damt2)
	_, C := test.RandKeys()
	b.dlc.oracleReqs.commitments[dID] = C

	tx, err := b.dlc.ContractExecutionTx(b.party, deal, dID)

	// asserions
	assert := assert.New(t)
	assert.Nil(tx)
	assert.NotNil(err)
	assert.IsType(&CETTakeNothingError{}, err)
}

func TestFixDeal(t *testing.T) {
	assert := assert.New(t)
	var err error

	// setup
	b, _, dID, deal := setupContractorsUntilPubkeyExchange(1, 1)

	// sucess with valid sign and message set
	privkey, C := test.RandKeys()
	b.dlc.oracleReqs.commitments[dID] = C
	osigns := [][]byte{privkey.D.Bytes()}
	osignset := &oracle.SignSet{Msgs: deal.Msgs, Signs: osigns}

	err = b.FixDeal(osignset)
	assert.NoError(err)

	// fail with invalid sign
	privkeyInvalid, _ := test.RandKeys()
	osignsInvalid := [][]byte{privkeyInvalid.D.Bytes()}
	osignsetInvalid := &oracle.SignSet{Msgs: deal.Msgs, Signs: osignsInvalid}

	err = b.FixDeal(osignsetInvalid)
	assert.Error(err)
}

func TestSignedContractExecutionTx(t *testing.T) {
	assert := assert.New(t)
	var err error

	// setup
	b1, b2, dID, deal := setupContractorsUntilPubkeyExchange(1, 1)
	privkey, C := test.RandKeys()
	b1.dlc.oracleReqs.commitments[dID] = C
	b2.dlc.oracleReqs.commitments[dID] = C
	osigns := [][]byte{privkey.D.Bytes()}
	osignset := &oracle.SignSet{Msgs: deal.Msgs, Signs: osigns}

	err = b1.FixDeal(osignset)
	assert.NoError(err)
	err = b2.FixDeal(osignset)
	assert.NoError(err)

	// fail without the counterparty's sign
	_, err = b1.SignedContractExecutionTx()
	assert.NoError(err)
	_, err = b2.SignedContractExecutionTx()
	assert.NoError(err)

	// exchange signs
	sign1, err := b1.SignContractExecutionTx(deal, dID)
	assert.NoError(err)
	sign2, err := b2.SignContractExecutionTx(deal, dID)
	assert.Nil(err)

	err = b1.AcceptCETxSigns([][]byte{sign2})
	assert.Nil(err)
	err = b2.AcceptCETxSigns([][]byte{sign1})
	assert.Nil(err)

	// no errors with the counterparty's sign
	tx1, err := b1.SignedContractExecutionTx()
	assert.NoError(err)
	tx2, err := b2.SignedContractExecutionTx()
	assert.NoError(err)

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

func setupContractorsUntilPubkeyExchange(
	damt1, damt2 btcutil.Amount) (b1, b2 *Builder, dID int, deal *Deal) {
	conds := newTestConditions()

	// set deals
	msgs := [][]byte{{1}}
	deal = NewDeal(damt1, damt2, [][]byte{{1}})
	conds.Deals = []*Deal{deal}

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

	dID, deal, _ = b1.dlc.DealByMsgs(msgs)

	return b1, b2, dID, deal
}

func runFundScript(b *Builder, tx *wire.MsgTx) error {
	d := b.DLC()
	fundtx, _ := d.FundTx()
	fout := fundtx.TxOut[fundTxOutAt]
	return test.ExecuteScript(fout.PkScript, tx, fout.Value)
}
