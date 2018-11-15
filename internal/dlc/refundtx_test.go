package dlc

import (
	"testing"

	"github.com/dgarage/dlc/internal/test"
	"github.com/stretchr/testify/assert"
)

const testLockTime = uint32(1541951794) // 2018/11/11 3:46pm (UTC)

func setupDLCRefund() (party1, party2 *Builder, d *DLC) {
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

	// sign refundtx
	rs1, _ := b1.SignRefundTx()
	rs2, _ := b2.SignRefundTx()

	// exchange refund signs
	b1.AcceptRefundTxSign(rs2)
	b2.AcceptRefundTxSign(rs1)

	d = b1.DLC()

	return b1, b2, d
}

// VerifyRefundTx should return false if given RefundTx isnt valid
func TestVerifyRefundTxBadRefund(t *testing.T) {
	assert := assert.New(t)

	_, _, d := setupDLCRefund()

	// VerifyRefundTX should return false bc the given signature and pubkey don't match
	testBadSign := []byte{'b', 'a', 'd'} // make a known bad signature
	err := d.VerifyRefundTx(testBadSign, d.pubs[FirstParty])
	assert.NotNil(err)
}

// VerifyRefundTx should return true if given RefundTx is valid
func TestVerifyRefundTx(t *testing.T) {
	assert := assert.New(t)

	_, _, d := setupDLCRefund()

	err := d.VerifyRefundTx(d.refundSigns[FirstParty], d.pubs[FirstParty])
	assert.Nil(err)
}

func TestRefundTx(t *testing.T) {
	assert := assert.New(t)

	_, _, d := setupDLCRefund()

	refundtx, err := d.RefundTx()
	assert.Nil(err)
	assert.Equal(refundtx.LockTime, testLockTime) // check lockTime is same as set by DLC
	assert.Len(refundtx.TxIn, 1)                  // fund from fundtx?
	assert.Len(refundtx.TxOut, 2)                 // 1 for party and 1 for counterparty

	// Both parties should be able to have their initial funds refunded.
	assert.Equal(refundtx.TxOut[0].Value, int64(d.conds.FundAmts[FirstParty]))
	assert.Equal(refundtx.TxOut[1].Value, int64(d.conds.FundAmts[SecondParty]))
}

// TODO: TestRedeemRefundTx? Test redeem before lock out time, test after?
func TestRefundTxOutput(t *testing.T) {
	assert := assert.New(t)
	_, _, d := setupDLCRefund()

	// run script
	refundtx, err := d.SignedRefundTx()
	assert.Nil(err)

	fundtx, _ := d.FundTx()
	fout := fundtx.TxOut[fundTxOutAt]
	err = test.ExecuteScript(fout.PkScript, refundtx, fout.Value)
	assert.Nil(err)
}
