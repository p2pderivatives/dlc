package dlc

import (
	"errors"

	"github.com/btcsuite/btcd/wire"
	"github.com/dgarage/dlc/internal/script"
)

// ContractExecutionTx constracts a contract execution tx (CET) using pubkeys and given condition.
// Both parties have different transactions signed by the other side.
//
// input:
//   [0]:fund transaction output[0]
// output:
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
	if amt1 > 0 {
		txout2, err := d.ClosingTxOut(cparty, amt2)
		if err != nil {
			return nil, err
		}
		tx.AddTxOut(txout2)
	}
	return tx, nil
}

// SignContractExecutionTx signs a contract exection tx for a given party
func (b *Builder) SignContractExecutionTx(
	party Contractor, idx int) ([]byte, error) {

	cparty := counterparty(b.party)

	deal, err := b.dlc.Deal(idx)
	if err != nil {
		return nil, err
	}

	tx, err := b.dlc.ContractExecutionTx(cparty, deal)
	if err != nil {
		return nil, err
	}

	// TODO: verify contract execution tx

	return b.witsigForFundTxIn(tx)
}

// SignedContractExecutionTx returns a contract execution tx signed by both parties
func (b *Builder) SignedContractExecutionTx(idx int) (*wire.MsgTx, error) {
	deal, err := b.dlc.Deal(idx)
	if err != nil {
		return nil, err
	}

	tx, err := b.dlc.ContractExecutionTx(b.party, deal)
	if err != nil {
		return nil, err
	}

	sign, err := b.witsigForFundTxIn(tx)
	if err != nil {
		return nil, err
	}

	if deal.cpSign == nil {
		return nil, errors.New("missing counterparty's sign")
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
