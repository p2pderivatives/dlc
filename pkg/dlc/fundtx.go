package dlc

import (
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/p2pderivatives/dlc/pkg/script"
)

const fundTxOutAt = 0 // fund txout is always at 0 in fund tx
const fundTxInAt = 0  // fund txin is always at 0 in redeem tx

// FundTx constructs fund tx using prepared fund tx requirements
func (d *DLC) FundTx() (*wire.MsgTx, error) {
	tx := wire.NewMsgTx(txVersion)

	txout, err := d.fundTxOutForRedeemTx()
	if err != nil {
		return nil, err
	}
	tx.AddTxOut(txout)

	for _, p := range []Contractor{FirstParty, SecondParty} {
		// txins
		total := btcutil.Amount(0)
		for _, utxo := range d.Utxos[p] {
			txin, err := utxoToTxIn(utxo)
			if err != nil {
				return nil, err
			}

			tx.AddTxIn(txin)
			amt, err := btcutil.NewAmount(utxo.Amount)
			if err != nil {
				return nil, err
			}
			total += amt
		}

		// txout for change
		change := total - d.DepositAmt(p)

		if change < 0 {
			msg := fmt.Sprintf("Not enough utxos from %s", p)
			return nil, errors.New(msg)
		}

		if change > 0 {
			addr := d.ChangeAddrs[p]
			if addr == nil {
				msg := fmt.Sprintf("Change address must be provided by %s", p)
				return nil, errors.New(msg)
			}
			sc, err := script.P2WPKHpkScriptFromAddress(addr)
			if err != nil {
				return nil, err
			}
			tx.AddTxOut(wire.NewTxOut(int64(change), sc))
		}
	}

	return tx, nil
}

func (d *DLC) fundScript() ([]byte, error) {
	pub1, ok := d.Pubs[FirstParty]
	if !ok {
		return nil, errors.New("First party must provide a pubkey for fund script")
	}
	pub2, ok := d.Pubs[SecondParty]
	if !ok {
		return nil, errors.New("Second party must provide a pubkey for fund script")
	}

	return script.FundScript(pub1, pub2)
}

// fundTxOutForRedeemTx creates a txout for the txin of redeem tx.
// The value of the txout is calculated by `fund amount + redeem tx fee`
func (d *DLC) fundTxOutForRedeemTx() (*wire.TxOut, error) {
	fs, err := d.fundScript()
	if err != nil {
		return nil, err
	}

	pkScript, err := script.P2WSHpkScript(fs)
	if err != nil {
		return nil, err
	}

	amt, err := d.fundAmount()
	if err != nil {
		return nil, err
	}

	amt += d.redeemTxFee(cetxSize)

	txout := wire.NewTxOut(int64(amt), pkScript)

	return txout, nil
}

func (d *DLC) witnessForFundScript(
	sig1, sig2 []byte) (wire.TxWitness, error) {

	sc, err := d.fundScript()
	if err != nil {
		return nil, err
	}

	wit := script.WitnessForFundScript(sig1, sig2, sc)
	return wit, nil
}

// fundAmount calculates total fund amount
func (d *DLC) fundAmount() (btcutil.Amount, error) {
	amt1, ok := d.Conds.FundAmts[FirstParty]
	if !ok {
		return 0, errors.New("Fund amount for first party isn't set")
	}
	amt2, ok := d.Conds.FundAmts[SecondParty]
	if !ok {
		return 0, errors.New("Fund amount for second party isn't set")
	}

	return amt1 + amt2, nil
}

// DepositAmt calculates fund amount + fees
func (d *DLC) DepositAmt(p Contractor) btcutil.Amount {
	famt := d.Conds.FundAmts[p]
	fee := d.totalFee(p)
	return famt + fee
}

// FundAmt returns fund amount
func (b *Builder) FundAmt() btcutil.Amount {
	return b.dlc.Conds.FundAmts[b.party]
}

// PrepareFundTx prepares fundtx ins and out
func (b *Builder) PrepareFundTx() error {
	famt := b.dlc.Conds.FundAmts[b.party]
	feeBase := b.dlc.fundTxFeeBase()
	redeemTxFee := b.dlc.redeemTxFee(cetxSize)
	utxos, change, err := b.wallet.SelectUnspent(
		famt+feeBase+redeemTxFee,
		b.dlc.fundTxFeePerTxIn(),
		b.dlc.fundTxFeePerTxOut())
	if err != nil {
		return err
	}

	// set utxos to DLC
	_utxos := []*Utxo{}
	for _, utxo := range utxos {
		_utxos = append(_utxos, &utxo)
	}
	b.dlc.Utxos[b.party] = _utxos

	if change > 0 {
		addr, err := b.wallet.NewAddress()
		if err != nil {
			return err
		}

		// set change addr to DLC
		b.dlc.ChangeAddrs[b.party] = addr
	}

	return nil
}

// Utxos returns utxos
func (b *Builder) Utxos() []Utxo {
	utxos := []Utxo{}
	for _, utxo := range b.dlc.Utxos[b.party] {
		utxos = append(utxos, *utxo)
	}
	return utxos
}

// AcceptUtxos accepts utxos
func (b *Builder) AcceptUtxos(utxos []Utxo) error {
	cp := counterparty(b.party)

	// TODO: validate if total amount is enough

	_utxos := []*Utxo{}
	for _, utxo := range utxos {
		_utxos = append(_utxos, &utxo)
	}
	b.dlc.Utxos[cp] = _utxos

	return nil
}

// ChangeAddress returns address to send change
func (b *Builder) ChangeAddress() btcutil.Address {
	addr := b.dlc.ChangeAddrs[b.party]
	return addr.(btcutil.Address)
}

// AcceptsChangeAdderss accepts change address from the counterparty
func (b *Builder) AcceptsChangeAdderss(addr btcutil.Address) {
	cp := counterparty(b.party)
	b.dlc.ChangeAddrs[cp] = addr
}

// newRedeemTx creates a new tx to redeem fundtx
// redeem tx
//  inputs:
//   [0]: fund transaction output[0]
func (d *DLC) newRedeemTx() (*wire.MsgTx, error) {
	fundtx, err := d.FundTx()
	if err != nil {
		return nil, err
	}
	// TODO: verify if fund tx is completed

	tx := wire.NewMsgTx(txVersion)

	// txin
	txid := fundtx.TxHash()
	fout := wire.NewOutPoint(&txid, fundTxOutAt)
	txin := wire.NewTxIn(fout, nil, nil)
	tx.AddTxIn(txin)

	return tx, nil
}

// witsigForFundScript returns signature for a given tx that redeems fund out
func (b *Builder) witsigForFundScript(tx *wire.MsgTx) ([]byte, error) {
	fundtx, err := b.dlc.FundTx()
	if err != nil {
		return nil, err
	}
	fout := fundtx.TxOut[fundTxOutAt]
	famt := btcutil.Amount(fout.Value)

	fc, err := b.dlc.fundScript()
	if err != nil {
		return nil, err
	}

	pub := b.dlc.Pubs[b.party]

	return b.wallet.WitnessSignature(tx, fundTxInAt, famt, fc, pub)
}

// SignFundTx signs fund tx and return witnesses for the txins owned by the party
func (b *Builder) SignFundTx() ([]wire.TxWitness, error) {
	fundtx, err := b.dlc.FundTx()
	if err != nil {
		return nil, err
	}

	// get witnesses
	idxs := b.dlc.fundTxInsIdxs(b.party)
	wits, err := b.wallet.WitnessSignTxByIdxs(fundtx, idxs)
	if err != nil {
		return nil, err
	}

	// set witnesses to dlc
	b.dlc.FundWits[b.party] = wits

	return wits, nil
}

// SignedFundTx constructs signed fundtx
func (d *DLC) SignedFundTx() (*wire.MsgTx, error) {
	tx, err := d.FundTx()
	if err != nil {
		return nil, err
	}

	for _, p := range []Contractor{FirstParty, SecondParty} {
		wits := d.FundWits[p]

		idxs := d.fundTxInsIdxs(p)

		if len(wits) != len(idxs) {
			msg := fmt.Sprintf(
				"Expected %d signatures from %s, but found %d", len(idxs), p, len(wits))
			return nil, errors.New(msg)
		}

		for i, idx := range idxs {
			tx.TxIn[idx].Witness = wits[i]
		}
	}

	return tx, nil
}

// SendFundTx sends fund tx to the network
func (b *Builder) SendFundTx() error {
	tx, err := b.dlc.SignedFundTx()
	if err != nil {
		return err
	}

	_, err = b.wallet.SendRawTransaction(tx)
	return err
}

// FundTxHex returns hex string of fund tx
func (b *Builder) FundTxHex() (string, error) {
	tx, err := b.dlc.SignedFundTx()
	if err != nil {
		return "", err
	}

	return txToHex(tx)
}

// fundTxInAt returns indices of txin in fundtx by the party
func (d *DLC) fundTxInsIdxs(p Contractor) (idxs []int) {
	nTxInMe := len(d.Utxos[p])
	var txinFrom, txinTo int
	if p == FirstParty {
		txinFrom = 0
		txinTo = nTxInMe
	} else {
		nTxInCP := len(d.Utxos[FirstParty])
		txinFrom = nTxInCP
		txinTo = txinFrom + nTxInMe
	}
	for i := txinFrom; i < txinTo; i++ {
		idxs = append(idxs, i)
	}
	return idxs
}

// AcceptFundWitnesses accepts witnesses for fund txins owned by the counerparty
func (b *Builder) AcceptFundWitnesses(wits []wire.TxWitness) {
	cparty := counterparty(b.party)
	b.dlc.FundWits[cparty] = wits

	// TODO: verify
}
