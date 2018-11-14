package dlc

import (
	"testing"

	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/dgarage/dlc/internal/mocks/walletmock"
	"github.com/dgarage/dlc/internal/test"
	"github.com/stretchr/testify/assert"
)

func setupDLC(conds Conditions) *DLC {
	d := newDLC(conds)
	_, pub1 := test.RandKeys()
	_, pub2 := test.RandKeys()
	d.pubs[FirstParty] = pub1
	d.pubs[SecondParty] = pub2
	return d
}

func newTestCETx(amt btcutil.Amount) *wire.MsgTx {
	tx := wire.NewMsgTx(txVersion)
	tx.AddTxOut(wire.NewTxOut(int64(amt), []byte{}))
	return tx
}

func TestClosingTxFailIfNotEnoughFees(t *testing.T) {
	conds, _ := NewConditions(1, 1, 1, 1, 1)
	d := setupDLC(conds)
	inamt := btcutil.Amount(1)
	cetx := newTestCETx(inamt)

	_, err := d.ClosingTx(FirstParty, cetx)

	assert := assert.New(t)
	assert.Error(err)
}

func TestClosingTx(t *testing.T) {
	conds, _ := NewConditions(1, 1, 1, 1, 1)
	d := setupDLC(conds)
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

func TestSignedClosingTx(t *testing.T) {
	assert := assert.New(t)

	// setup
	b1, b2 := setupContractorsUntilSignExchange()

	// first party
	deal1, _ := b1.dlc.FixedDeal()
	cetx1, _ := b1.SignedContractExecutionTx()
	tx1, err := b1.SignedClosingTx(deal1, cetx1)
	assert.NoError(err)

	// second party
	deal2, _ := b2.dlc.FixedDeal()
	cetx2, _ := b2.SignedContractExecutionTx()
	tx2, err := b2.SignedClosingTx(deal2, cetx2)
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

func setupContractorsUntilSignExchange() (b1, b2 *Builder) {
	conds, _ := NewConditions(1, 1, 1, 1, 1)
	msgPriv, msgCommit := test.RandKeys()
	msgSign := msgPriv.D.Bytes()

	// init first party
	w1 := setupTestWalletForTestSignedClosingTx(msgSign)
	b1 = NewBuilder(FirstParty, w1, conds)
	b1.PreparePubkey()
	b1.PrepareFundTxIns()

	// init second party
	w2 := setupTestWalletForTestSignedClosingTx(msgSign)
	b2 = NewBuilder(SecondParty, w2, conds)
	b2.PreparePubkey()
	b2.PrepareFundTxIns()

	// exchange pubkeys
	b1.CopyReqsFromCounterparty(b2.DLC())
	b2.CopyReqsFromCounterparty(b1.DLC())

	dID, deal := setupDeal(b1,
		1*btcutil.SatoshiPerBitcoin, 1*btcutil.SatoshiPerBitcoin)
	b1.SetMsgCommitmentToDeal(dID, msgCommit)
	_, _ = setupDeal(b2,
		1*btcutil.SatoshiPerBitcoin, 1*btcutil.SatoshiPerBitcoin)
	b2.SetMsgCommitmentToDeal(dID, msgCommit)

	msign := msgPriv.D.Bytes()
	b1.FixDeal(dID, msign)
	b2.FixDeal(dID, msign)

	sign1, _ := b1.SignContractExecutionTx(deal)
	sign2, _ := b2.SignContractExecutionTx(deal)

	_ = b1.AcceptCETxSign(dID, sign2)
	_ = b2.AcceptCETxSign(dID, sign1)

	return b1, b2
}

// setup mocke wallet
func setupTestWalletForTestSignedClosingTx(msgSign []byte) *walletmock.Wallet {
	w := &walletmock.Wallet{}
	priv, pub := test.RandKeys()
	w.On("NewPubkey").Return(pub, nil)
	w = mockWitnessSignature(w, pub, priv)
	w = mockSelectUnspent(w, 1, 1, nil)
	w = mockWitnessSignatureWithCallback(
		w, pub, priv, genAddSignToPrivkeyFunc(msgSign))
	return w
}

func runCEScript(cetx *wire.MsgTx, tx *wire.MsgTx) error {
	cetxout := cetx.TxOut[closingTxOutAt]
	return test.ExecuteScript(cetxout.PkScript, tx, cetxout.Value)
}