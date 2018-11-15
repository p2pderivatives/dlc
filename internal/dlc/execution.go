package dlc

import (
	"errors"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/dgarage/dlc/internal/script"
)

// ContractExecutionTx constructs a contract execution tx (CET) using pubkeys and given condition.
// Both parties have different transactions signed by the other side.
//
// txins:
//   [0]:fund transaction output[0]
// txouts:
//   [0]:settlement script
//   [1]:p2wpkh (option)
func (d *DLC) ContractExecutionTx(
	party Contractor, deal *Deal) (*wire.MsgTx, error) {
	cparty := counterparty(party)

	tx, err := d.newRedeemTx()
	if err != nil {
		return nil, err
	}

	// out values
	amt1 := deal.amts[party]
	amt2 := deal.amts[cparty]

	if amt1 == 0 {
		errmsg := "Amount for a multisig script address shouldn't be zero"
		return nil, newCETTakeNothingError(errmsg)
	}

	// txout1: contract execution script
	pub1 := d.pubs[party]
	if pub1 == nil {
		return nil, errors.New("missing pubkey")
	}
	pub2 := d.pubs[cparty]
	if pub2 == nil {
		return nil, errors.New("missing pubkey")
	}
	if deal.msgCommitment == nil {
		return nil, errors.New("missing oracle's message commitment")
	}
	sc, err := script.ContractExecutionScript(pub1, pub2, deal.msgCommitment)
	if err != nil {
		return nil, err
	}
	pkScript, err := script.P2WSHpkScript(sc)
	if err != nil {
		return nil, err
	}

	txout1 := wire.NewTxOut(int64(amt1), pkScript)
	tx.AddTxOut(txout1)

	// txout2: counterparty's p2wpkh
	if amt2 > 0 {
		txout2, err := d.ClosingTxOut(cparty, amt2)
		if err != nil {
			return nil, err
		}
		tx.AddTxOut(txout2)
	}
	return tx, nil
}

// SignContractExecutionTx signs a contract execution tx for a given party
func (b *Builder) SignContractExecutionTx(deal *Deal) ([]byte, error) {
	cparty := counterparty(b.party)

	tx, err := b.dlc.ContractExecutionTx(cparty, deal)
	if err != nil {
		return nil, err
	}

	return b.witsigForFundTxIn(tx)
}

// AcceptCETxSign sets a sign received from the counterparty
func (b *Builder) AcceptCETxSign(
	idx int, sign []byte) error {

	d, err := b.dlc.Deal(idx)
	if err != nil {
		return err
	}

	// verify
	ok, err := b.dlc.verifyContractExecutionSign(b.party, d, sign)
	if !ok {
		return err
	}

	d.cpSign = sign
	return nil
}

func (d *DLC) verifyContractExecutionSign(
	p Contractor, deal *Deal, sign []byte) (bool, error) {

	tx, err := d.ContractExecutionTx(p, deal)
	if err != nil {
		return false, err
	}

	cparty := counterparty(p)

	fsc, err := d.fundScript()
	if err != nil {
		return false, err
	}

	sighashes := txscript.NewTxSigHashes(tx)

	ftx, err := d.FundTx()
	if err != nil {
		return false, err
	}

	fout := ftx.TxOut[fundTxOutAt]

	hash, err := txscript.CalcWitnessSigHash(
		fsc, sighashes, txscript.SigHashAll, tx, fundTxInAt, fout.Value)
	if err != nil {
		return false, err
	}

	s, err := btcec.ParseDERSignature(sign, btcec.S256())
	if err != nil {
		return false, err
	}

	if !s.Verify(hash, d.pubs[cparty]) {
		return false, errors.New("failed to verify")
	}

	return true, nil
}

// SignedContractExecutionTx returns a contract execution tx signed by both parties
func (b *Builder) SignedContractExecutionTx() (*wire.MsgTx, error) {
	deal, err := b.dlc.FixedDeal()
	if err != nil {
		return nil, err
	}

	tx, err := b.dlc.ContractExecutionTx(b.party, deal)
	if err != nil {
		return nil, err
	}

	if deal.cpSign == nil {
		return nil, errors.New("missing counterparty's sign")
	}

	sign, err := b.witsigForFundTxIn(tx)
	if err != nil {
		return nil, err
	}

	var sign1, sign2 []byte
	switch b.party {
	case FirstParty:
		sign1, sign2 = sign, deal.cpSign
	case SecondParty:
		sign1, sign2 = deal.cpSign, sign
	}

	wit, err := b.dlc.witnessForFundScript(sign1, sign2)
	if err != nil {
		return nil, err
	}

	tx.TxIn[fundTxInAt].Witness = wit

	return tx, nil
}
