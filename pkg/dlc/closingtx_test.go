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

func setupContractorsUntilSignExchange() (b1, b2 *Builder) {
	conds := newTestConditions()

	var damt1, damt2 btcutil.Amount = 1 * btcutil.SatoshiPerBitcoin, 1 * btcutil.SatoshiPerBitcoin
	msgs := [][]byte{{1}}
	deal := NewDeal(damt1, damt2, msgs)
	conds.Deals = []*Deal{deal}

	// oracle's signnature and commitment
	opriv, C := test.RandKeys()
	osig := opriv.D.Bytes()
	oFixedMsg := &oracle.SignedMsg{Msgs: msgs, Sigs: [][]byte{osig}}

	// init first party
	w1 := setupTestWalletForTestSignedClosingTx(osig)
	b1 = NewBuilder(FirstParty, w1, conds)
	b1.PreparePubkey()
	b1.PrepareFundTx()

	// init second party
	w2 := setupTestWalletForTestSignedClosingTx(osig)
	b2 = NewBuilder(SecondParty, w2, conds)
	b2.PreparePubkey()
	b2.PrepareFundTx()

	// exchange pubkeys and utxos
	stepSendRequirments(b1, b2)
	stepSendRequirments(b2, b1)

	dID, _, _ := b1.Contract.DealByMsgs(msgs)

	// set oracle ocmmitment
	b1.Contract.Oracle.Commitments[dID] = C
	b2.Contract.Oracle.Commitments[dID] = C

	// exchange sigs
	sig1, _ := b1.SignContractExecutionTx(deal, dID)
	sig2, _ := b2.SignContractExecutionTx(deal, dID)
	_ = b1.AcceptCETxSignatures([][]byte{sig2})
	_ = b2.AcceptCETxSignatures([][]byte{sig1})

	// fix deal by oracle's sig
	b1.FixDeal(oFixedMsg, []int{0})
	b2.FixDeal(oFixedMsg, []int{0})

	return b1, b2
}

// setup mocke wallet
func setupTestWalletForTestSignedClosingTx(msgSig []byte) *walletmock.Wallet {
	w := &walletmock.Wallet{}
	w = mockNewAddress(w)
	w = mockSelectUnspent(w, 1000, 1, nil)

	priv, pub := test.RandKeys()
	w.On("NewPubkey").Return(pub, nil)
	w = mockWitnessSignature(w, pub, priv)
	w = mockWitnessSignatureWithCallback(
		w, pub, priv, genAddSigToPrivkeyFunc(msgSig))

	return w
}

func runCEScript(cetx *wire.MsgTx, tx *wire.MsgTx) error {
	cetxout := cetx.TxOut[closingTxOutAt]
	return test.ExecuteScript(cetxout.PkScript, tx, cetxout.Value)
}
