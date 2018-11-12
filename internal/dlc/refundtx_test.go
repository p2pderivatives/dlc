package dlc

import (
	"fmt"
	"testing"

	"github.com/btcsuite/btcd/wire"
	"github.com/dgarage/dlc/internal/test"
	"github.com/stretchr/testify/assert"
)

var (
	testLockTime = uint32(1541951794) // this is equivalent to 2018/11/11 3:46pm (UTC)
	// TODO: add block # testlocktime?
)

func setupDLCRefund() (party1, party2 *Builder, d *DLC) {
	// init first party
	w1 := setupTestWallet()
	b1 := NewBuilder(FirstParty, mockSelectUnspent(w1, 1, 1, nil))
	b1.SetFundAmounts(1, 1)
	b1.PreparePubkey()
	b1.PrepareFundTxIns()

	// init second party
	w2 := setupTestWallet()
	b2 := NewBuilder(SecondParty, mockSelectUnspent(w2, 1, 1, nil))
	b2.SetFundAmounts(1, 1)
	b2.PreparePubkey()
	b2.PrepareFundTxIns()

	// set locktime
	b2.SetLockTime(testLockTime) // below locktime threshold 5e8
	// b1.SetLockTime(uint32(8e8)) // below locktime threshold

	// exchange pubkeys
	b1.CopyReqsFromCounterparty(b2.DLC())
	b2.CopyReqsFromCounterparty(b1.DLC())

	// sign refundtx
	b2.SignRefundTx()
	b1.SignRefundTx()

	// exchange refund signs
	b1.CopyReqsFromCounterparty(b2.DLC())
	b2.CopyReqsFromCounterparty(b1.DLC())

	d = b1.DLC()

	return b1, b2, d
}

// VerifyRefundTx should return false if given RefundTx isnt valid
func TestVerifyRefundTxBadRefund(t *testing.T) {
	assert := assert.New(t)

	_, _, d := setupDLCRefund()

	valid, err := d.VerifyRefundTx(d.refundSigns[SecondParty], d.pubs[FirstParty])
	assert.NotNil(err)
	assert.False(valid)
}

// VerifyRefundTx should return true if given RefundTx is valid
func TestVerifyRefundTx(t *testing.T) {
	assert := assert.New(t)

	_, _, d := setupDLCRefund()

	valid, err := d.VerifyRefundTx(d.refundSigns[FirstParty], d.pubs[FirstParty])
	assert.Nil(err)
	assert.True(valid)
}

func TestRefundTx(t *testing.T) {
	assert := assert.New(t)

	_, _, d := setupDLCRefund()

	refundtx, err := d.RefundTx()
	assert.Nil(err)
	assert.Equal(refundtx.LockTime, testLockTime) // check lockTime is same as set by DLC
	assert.Len(refundtx.TxIn, 1)                  // fund from fundtx?
	assert.Len(refundtx.TxOut, 2)                 // 1 for party and 1 for counterparty
}

// TestRedeemRefundTx? Test redeem before lock out time, test after?
func TestRedeemRefundTx(t *testing.T) {
	assert := assert.New(t)
	b1, b2, d := setupDLCRefund()

	// prepare redeem tx for testing. this will be a settlement tx or refund tx
	redeemtx, err := d.newRedeemTx() // CHANGE THIS
	assert.Nil(err)
	// both parties signs redeem tx
	sign1, err := b1.witsigForFundTxIn(redeemtx)
	assert.Nil(err)
	sign2, err := b2.witsigForFundTxIn(redeemtx)
	assert.Nil(err)

	// create witness
	fsc, _ := d.fundScript()
	wt := wire.TxWitness{[]byte{}, sign1, sign2, fsc}
	redeemtx.TxIn[0].Witness = wt

	// run script
	refundtx, _ := d.RefundTx()
	fmt.Printf("REFUNDTX OUT:    %+v\n", refundtx.TxOut)

	fout := refundtx.TxOut[1]
	err = test.ExecuteScript(fout.PkScript, redeemtx, fout.Value)
	assert.Nil(err)

}
