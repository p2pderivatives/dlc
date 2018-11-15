package dlc

import (
	"testing"

	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/dgarage/dlc/internal/mocks/walletmock"
	"github.com/dgarage/dlc/internal/test"
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
	d := newDLC(newTestConditions())
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

func TestSignedClosingTx(t *testing.T) {
	assert := assert.New(t)

	// setup
	b1, b2 := setupContractorsUntilSignExchange()

	// first party
	cetx1, _ := b1.SignedContractExecutionTx()
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

func setupContractorsUntilSignExchange() (b1, b2 *Builder) {
	conds := newTestConditions()

	var damt1, damt2 btcutil.Amount = 1 * btcutil.SatoshiPerBitcoin, 1 * btcutil.SatoshiPerBitcoin
	msgs := [][]byte{{1}}
	deal := NewDeal(damt1, damt2, msgs)
	conds.Deals = []*Deal{deal}

	msgPriv, _ := test.RandKeys()
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

	dID, _, _ := b1.dlc.DealByMsgs(msgs)
	// TODO: fix
	// b1.FixDeal(dID, msign)
	// b2.FixDeal(dID, msign)

	sign1, _ := b1.SignContractExecutionTx(deal, dID)
	sign2, _ := b2.SignContractExecutionTx(deal, dID)

	_ = b1.AcceptCETxSigns([][]byte{sign2})
	_ = b2.AcceptCETxSigns([][]byte{sign1})

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
