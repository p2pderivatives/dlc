package dlc

import (
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/dgarage/dlc/internal/script"
)

// ClosingTxSize is size of closing tx
const ClosingTxSize = 216

// ClosingTxOutAt is a txout index of contract execution tx
const ClosingTxOutAt = 0

// ClosingTx constructs a tx that redeems a given CET
func (d *DLC) ClosingTx(
	p Contractor, cetx *wire.MsgTx) (*wire.MsgTx, error) {

	tx := wire.NewMsgTx(txVersion)

	// txin
	txid := cetx.TxHash()
	txin := wire.NewTxIn(
		wire.NewOutPoint(&txid, ClosingTxOutAt), nil, nil)

	tx.AddTxIn(txin)

	// TODO: the party who sends closing tx have to pay for this fee.
	// Consider spliting it with the counterparty
	amt := d.fundAmts[p] - d.closingTxFee()
	txout, err := d.ClosingTxOut(p, amt)
	if err != nil {
		return nil, err
	}
	tx.AddTxOut(txout)

	return tx, nil
}

func (d *DLC) closingTxFee() btcutil.Amount {
	return d.redeemFeerate.MulF64(ClosingTxSize)
}

// SignedClosingTx constructs a closing tx with witness
func (b *Builder) SignedClosingTx() (*wire.MsgTx, error) {
	deal, err := b.dlc.FixedDeal()
	if err != nil {
		return nil, err
	}

	cetx, err := b.dlc.ContractExecutionTx(b.party, deal)
	if err != nil {
		return nil, err
	}

	tx, err := b.dlc.ClosingTx(b.party, cetx)
	if err != nil {
		return nil, err
	}

	wit, err := b.witnessForCEScript(tx, cetx, deal)
	if err != nil {
		return nil, err
	}
	tx.TxIn[0].Witness = wit

	return tx, nil
}

func (b *Builder) witnessForCEScript(tx, cetx *wire.MsgTx, deal *Deal) (wire.TxWitness, error) {
	cetxout := cetx.TxOut[ClosingTxOutAt]
	amt := cetxout.Value

	cparty := counterparty(b.party)
	pub1, pub2 := b.dlc.pubs[b.party], b.dlc.pubs[cparty]
	sc, err := script.ContractExecutionScript(
		pub1, pub2, deal.msgCommitment)
	if err != nil {
		return nil, err
	}

	privkeyConverter := func(priv *btcec.PrivateKey) (*btcec.PrivateKey, error) {
		// TODO: add msg sign
		return priv, nil
	}

	sign, err := b.wallet.WitnessSignatureWithCallback(
		tx, ClosingTxOutAt, amt, sc, pub1, privkeyConverter)
	if err != nil {
		return nil, err
	}

	wit := script.WitnessForCEScript(sign, sc)
	return wit, nil
}
