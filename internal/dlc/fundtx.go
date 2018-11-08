package dlc

import (
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/dgarage/dlc/internal/script"
	"github.com/dgarage/dlc/internal/wallet"
)

// FundTxRequirements contains txins and txouts for fund tx
type FundTxRequirements struct {
	txIns map[Contractor][]*wire.TxIn
	txOut map[Contractor]*wire.TxOut
	pubs  map[Contractor]*btcec.PublicKey
}

func newFundTxReqs() *FundTxRequirements {
	return &FundTxRequirements{
		txIns: make(map[Contractor][]*wire.TxIn),
		txOut: make(map[Contractor]*wire.TxOut),
		pubs:  make(map[Contractor]*btcec.PublicKey),
	}
}

// CopyFundTxReqsFromCounterparty copies requirements of counterparty from their DLC
func (b *Builder) CopyFundTxReqsFromCounterparty(d *DLC) {
	p := b.counterparty()
	reqs := d.fundTxReqs

	b.dlc.fundTxReqs.txIns[p] = reqs.txIns[p]
	b.dlc.fundTxReqs.txOut[p] = reqs.txOut[p]
	b.dlc.fundTxReqs.pubs[p] = reqs.pubs[p]
}

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
		for _, txin := range d.fundTxReqs.txIns[p] {
			tx.AddTxIn(txin)
		}
		// txout for change
		txout := d.fundTxReqs.txOut[p]
		if txout != nil {
			tx.AddTxOut(txout)
		}
	}

	return tx, nil
}

func (d *DLC) fundScript() ([]byte, error) {
	pub1, ok := d.fundTxReqs.pubs[FirstParty]
	if !ok {
		return nil, errors.New("First party must provide a pubkey for fund script")
	}
	pub2, ok := d.fundTxReqs.pubs[SecondParty]
	if !ok {
		return nil, errors.New("Second party must provide a pubkey for fund script")
	}

	return script.MultiSigScript2of2(pub1, pub2)
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

	amt += d.redeemTxFee()

	txout := wire.NewTxOut(int64(amt), pkScript)

	return txout, nil
}

// SetFundAmounts sets fund amounts to DLC
func (b *Builder) SetFundAmounts(amt1, amt2 btcutil.Amount) {
	b.dlc.fundAmts[FirstParty] = amt1
	b.dlc.fundAmts[SecondParty] = amt2
}

// fundAmount calculates total fund amount
func (d *DLC) fundAmount() (btcutil.Amount, error) {
	amt1, ok := d.fundAmts[FirstParty]
	if !ok {
		return 0, errors.New("Fund amount for first party isn't set")
	}
	amt2, ok := d.fundAmts[SecondParty]
	if !ok {
		return 0, errors.New("Fund amount for second party isn't set")
	}

	return amt1 + amt2, nil
}

// SetFundFeerate sets feerate (satoshi/byte) for fund tx fee calculation
func (b *Builder) SetFundFeerate(feerate btcutil.Amount) {
	b.dlc.fundFeerate = feerate
}

// SetRedeemFeerate sets feerate (satoshi/byte) for fund tx fee calculation
func (b *Builder) SetRedeemFeerate(feerate btcutil.Amount) {
	b.dlc.redeemFeerate = feerate
}

// Tx sizes for fee estimation
const fundTxBaseSize = int64(55)
const fundTxInSize = int64(149)
const fundTxOutSize = int64(31)
const redeemTxSize = int64(345)

func (d *DLC) fundTxFeeBase() btcutil.Amount {
	return d.fundFeerate.MulF64(float64(fundTxBaseSize))
}

func (d *DLC) fundTxFeePerTxIn() btcutil.Amount {
	return d.fundFeerate.MulF64(float64(fundTxInSize))
}

func (d *DLC) fundTxFeePerTxOut() btcutil.Amount {
	return d.fundFeerate.MulF64(float64(fundTxOutSize))
}

func (d *DLC) redeemTxFee() btcutil.Amount {
	return d.redeemFeerate.MulF64(float64(redeemTxSize))
}

// PrepareFundTxIns prepares utxos for fund tx by calculating fees
func (b *Builder) PrepareFundTxIns() error {
	famt, ok := b.dlc.fundAmts[b.party]
	if !ok {
		err := fmt.Errorf("fund amount isn't set yet")
		return err
	}

	feeBase := b.dlc.fundTxFeeBase()

	// TODO: add redeem tx fee

	feePerIn := b.dlc.fundTxFeePerTxIn()
	feePerOut := b.dlc.fundTxFeePerTxOut()
	utxos, change, err := b.wallet.SelectUnspent(famt+feeBase, feePerIn, feePerOut)
	if err != nil {
		return err
	}

	txins, err := wallet.UtxosToTxIns(utxos)
	if err != nil {
		return err
	}

	// set txins to DLC
	b.dlc.fundTxReqs.txIns[b.party] = txins

	if change > 0 {
		pub, err := b.wallet.NewPubkey()
		if err != nil {
			return err
		}

		// TODO: manager pubkey address for change

		pkScript, err := script.P2WPKHpkScript(pub)
		if err != nil {
			return err
		}

		txout := wire.NewTxOut(int64(change), pkScript)

		// set change txout to DLC
		b.dlc.fundTxReqs.txOut[b.party] = txout
	}

	return nil
}

// PrepareFundPubkey sets fund pubkey
func (b *Builder) PrepareFundPubkey() error {
	pub, err := b.wallet.NewPubkey()
	if err != nil {
		return err
	}
	b.dlc.fundTxReqs.pubs[b.party] = pub
	return nil
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

	tx := wire.NewMsgTx(txVersion)

	// txin
	txid := fundtx.TxHash()
	fout := wire.NewOutPoint(&txid, fundTxOutAt)
	txin := wire.NewTxIn(fout, nil, nil)
	tx.AddTxIn(txin)

	return tx, nil
}

// witsigForFundTxIn returns sign for a given tx that redeems fund out
// TODO: this method will be used to create settlement txs and refund tx
func (b *Builder) witsigForFundTxIn(tx *wire.MsgTx) ([]byte, error) {
	amt, err := b.dlc.fundAmount()
	if err != nil {
		return nil, err
	}

	fc, err := b.dlc.fundScript()
	if err != nil {
		return nil, err
	}

	pub := b.dlc.fundTxReqs.pubs[b.party]

	return b.wallet.WitnessSignature(tx, fundTxInAt, amt, fc, pub)
}
