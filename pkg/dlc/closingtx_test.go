package dlc

import (
	"testing"

	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/p2pderivatives/dlc/internal/mocks/walletmock"
	"github.com/p2pderivatives/dlc/internal/oracle"
	"github.com/p2pderivatives/dlc/internal/test"
	"github.com/stretchr/testify/assert"
)

func TestClosingTxFailIfNotEnoughFees(t *testing.T) {
	d := setupDLC()
	inamt := btcutil.Amount(1)
	cetx := newTestCETx(inamt)

	_, err := d.ClosingTx(FirstParty, cetx)

	assert := assert.New(t)
	assert.Error(err)
}

func TestClosingTx(t *testing.T) {
	d := setupDLC()
	inamt := btcutil.Amount(1 * btcutil.SatoshiPerBitcoin)
	cetx := newTestCETx(inamt)

	tx, err := d.ClosingTx(FirstParty, cetx)

	assert := assert.New(t)
	assert.NoError(err)
	assert.Len(tx.TxOut, 1)
	assert.InDelta(
		cetx.TxOut[closingTxOutAt].Value,
		tx.TxOut[0].Value,
		100000, // satoshi
	)
}

func setupDLC() *DLC {
	d := NewDLC(newTestConditions())
	_, pub1 := test.RandKeys()
	_, pub2 := test.RandKeys()
	d.Pubs[FirstParty] = pub1
	d.Pubs[SecondParty] = pub2
	d.Addrs[FirstParty] = test.RandAddress()
	d.Addrs[SecondParty] = test.RandAddress()
	return d
}

func newTestCETx(amt btcutil.Amount) *wire.MsgTx {
	tx := wire.NewMsgTx(txVersion)
	tx.AddTxOut(wire.NewTxOut(int64(amt), []byte{}))
	return tx
}

func TestSignedClosingTx(t *testing.T) {
	assert := assert.New(t)

	// setup
	b1, b2, err := setupContractorsUntilSignExchange()
	if !assert.NoError(err) {
		assert.FailNow(err.Error())
	}

	// first party
	cetx1, err := b1.SignedContractExecutionTx()
	assert.NoError(err)
	assert.NotEmpty(cetx1)
	tx1, err := b1.SignedClosingTx(cetx1)
	assert.NoError(err)

	// second party
	cetx2, _ := b2.SignedContractExecutionTx()
	tx2, err := b2.SignedClosingTx(cetx2)
	assert.NoError(err)

	// first party can redeem only their tx
	err = runCEScript(cetx1, tx1)
	assert.NoError(err)
	err = runCEScript(cetx2, tx1)
	assert.Error(err)

	// second party can redeem only their tx
	err = runCEScript(cetx2, tx2)
	assert.NoError(err)
	err = runCEScript(cetx1, tx2)
	assert.Error(err)
}

func setupContractorsUntilSignExchange() (b1, b2 *Builder, err error) {
	// msg
	msgs := [][]byte{{1}}
	damt1 := btcutil.Amount(1 * btcutil.SatoshiPerBitcoin)
	damt2 := btcutil.Amount(1 * btcutil.SatoshiPerBitcoin)
	deal := NewDeal(damt1, damt2, msgs)

	// oracle's signnature and commitment
	opriv, C := test.RandKeys()
	osig := opriv.D.Bytes()
	oFixedMsg := &oracle.SignedMsg{
		Msgs: msgs, Sigs: [][]byte{osig}}

	setupConds := func() *Conditions {
		conds := newTestConditions()
		conds.Deals = []*Deal{deal}
		return conds
	}

	setupWallet := func() *walletmock.Wallet {
		w := &walletmock.Wallet{}

		priv, pub := test.RandKeys()
		w.On("NewPubkey").Return(pub, nil)
		w = mockWitnessSignature(w, pub, priv)
		w = mockWitnessSignatureWithCallback(
			w, pub, priv, genAddSigToPrivkeyFunc(osig))

		return w
	}

	// init first party
	b1 = setupBuilder(FirstParty, setupWallet, setupConds)
	b2 = setupBuilder(SecondParty, setupWallet, setupConds)

	// preps
	if err = stepPrepare(b1); err != nil {
		return
	}
	if err = stepPrepare(b2); err != nil {
		return
	}

	// exchange pubkeys and utxos
	if err = stepSendRequirments(b1, b2); err != nil {
		return
	}
	if err = stepSendRequirments(b2, b1); err != nil {
		return
	}

	dID, _, err := b1.Contract.DealByMsgs(msgs)
	if err != nil {
		return
	}

	// set oracle ocmmitment
	b1.Contract.Oracle.Commitments[dID] = C
	b2.Contract.Oracle.Commitments[dID] = C

	// exchange sigs
	if err = stepExchangeCETxSig(b1, b2, deal, dID); err != nil {
		return
	}
	if err = stepExchangeCETxSig(b2, b1, deal, dID); err != nil {
		return
	}

	// fix deal by oracle's sig
	if err = b1.FixDeal(oFixedMsg, []int{0}); err != nil {
		return
	}
	if err = b2.FixDeal(oFixedMsg, []int{0}); err != nil {
		return
	}

	return b1, b2, nil
}

func runCEScript(cetx *wire.MsgTx, tx *wire.MsgTx) error {
	cetxout := cetx.TxOut[closingTxOutAt]
	return test.ExecuteScript(cetxout.PkScript, tx, cetxout.Value)
}
