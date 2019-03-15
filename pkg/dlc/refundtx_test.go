package dlc

import (
	"testing"

	"github.com/p2pderivatives/dlc/internal/test"
	"github.com/stretchr/testify/assert"
)

const testLockTime = uint32(1541951794) // 2018/11/11 3:46pm (UTC)

func setupDLCRefund() (b1, b2 *Builder, d *DLC, err error) {
	setupConds := func() *Conditions {
		conds := newTestConditions()
		conds.RefundLockTime = testLockTime
		return conds
	}

	// init builders
	b1 = setupBuilder(FirstParty, setupTestWallet, setupConds)
	b2 = setupBuilder(SecondParty, setupTestWallet, setupConds)

	if err = stepPrepare(b1); err != nil {
		return
	}
	if err = stepPrepare(b2); err != nil {
		return
	}

	// exchange pubkeys
	if err = stepSendRequirments(b1, b2); err != nil {
		return
	}

	if err = stepSendRequirments(b2, b1); err != nil {
		return
	}

	// sign refundtx
	rs1, err := b1.SignRefundTx()
	if err != nil {
		return
	}
	rs2, err := b2.SignRefundTx()
	if err != nil {
		return
	}

	// exchange refund signs
	if err = b1.AcceptRefundTxSignature(rs2); err != nil {
		return
	}

	if err = b2.AcceptRefundTxSignature(rs1); err != nil {
		return
	}

	return b1, b2, b1.Contract, nil
}

// VerifyRefundTx should return false if given RefundTx isnt valid
func TestVerifyRefundTxInvalidSig(t *testing.T) {
	assert := assert.New(t)

	_, _, d, err := setupDLCRefund()
	if !assert.NoError(err) {
		assert.FailNow(err.Error())
	}

	// VerifyRefundTX should return false bc the given signature and pubkey don't match
	testBadSig := []byte{'b', 'a', 'd'} // make a known bad signature
	err = d.VerifyRefundTx(testBadSig, d.Pubs[FirstParty])
	assert.NotNil(err)
}

// VerifyRefundTx should return true if given RefundTx is valid
func TestVerifyRefundTx(t *testing.T) {
	assert := assert.New(t)

	_, _, d, err := setupDLCRefund()
	if !assert.NoError(err) {
		assert.FailNow(err.Error())
	}

	err = d.VerifyRefundTx(d.RefundSigs[FirstParty], d.Pubs[FirstParty])
	assert.Nil(err)
}

func TestRefundTx(t *testing.T) {
	assert := assert.New(t)

	_, _, d, err := setupDLCRefund()
	if !assert.NoError(err) {
		assert.FailNow(err.Error())
	}

	refundtx, err := d.RefundTx()
	assert.Nil(err)
	assert.Equal(testLockTime, refundtx.LockTime) // check lockTime is same as set by DLC
	assert.Len(refundtx.TxIn, 1)                  // fund from fundtx
	assert.Len(refundtx.TxOut, 2)                 // 1 for party and 1 for counterparty

	// Both parties should be able to have their initial funds refunded.
	assert.Equal(refundtx.TxOut[0].Value, int64(d.Conds.FundAmts[FirstParty]))
	assert.Equal(refundtx.TxOut[1].Value, int64(d.Conds.FundAmts[SecondParty]))
}

func TestRefundTxOutput(t *testing.T) {
	assert := assert.New(t)
	_, _, d, err := setupDLCRefund()
	if !assert.NoError(err) {
		assert.FailNow(err.Error())
	}

	// run script
	refundtx, err := d.SignedRefundTx()
	assert.Nil(err)

	fundtx, _ := d.FundTx()
	fout := fundtx.TxOut[fundTxOutAt]
	err = test.ExecuteScript(fout.PkScript, refundtx, fout.Value)
	assert.Nil(err)
}
