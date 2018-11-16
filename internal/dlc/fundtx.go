package dlc

import (
	"errors"

	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/dgarage/dlc/internal/script"
	"github.com/dgarage/dlc/internal/wallet"
)

// FundTxRequirements contains txins and txouts for fund tx
type FundTxRequirements struct {
	txIns map[Contractor][]*wire.TxIn
	txOut map[Contractor]*wire.TxOut
}

func newFundTxReqs() *FundTxRequirements {
	return &FundTxRequirements{
		txIns: make(map[Contractor][]*wire.TxIn),
		txOut: make(map[Contractor]*wire.TxOut),
	}
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
	pub1, ok := d.pubs[FirstParty]
	if !ok {
		return nil, errors.New("First party must provide a pubkey for fund script")
	}
	pub2, ok := d.pubs[SecondParty]
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
	sign1, sign2 []byte) (wire.TxWitness, error) {

	sc, err := d.fundScript()
	if err != nil {
		return nil, err
	}

	wit := script.WitnessForFundScript(sign1, sign2, sc)
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

// Tx sizes for fee estimation
const fundTxBaseSize = int64(55)
const fundTxInSize = int64(149)
const fundTxOutSize = int64(31)
const cetxSize = int64(345) // context execution tx size

func (d *DLC) fundTxFeeBase() btcutil.Amount {
	return d.Conds.FundFeerate.MulF64(float64(fundTxBaseSize))
}

func (d *DLC) fundTxFeePerTxIn() btcutil.Amount {
	return d.Conds.FundFeerate.MulF64(float64(fundTxInSize))
}

func (d *DLC) fundTxFeePerTxOut() btcutil.Amount {
	return d.Conds.FundFeerate.MulF64(float64(fundTxOutSize))
}

func (d *DLC) redeemTxFee(size int64) btcutil.Amount {
	return d.Conds.RedeemFeerate.MulF64(float64(size))
}

// PrepareFundTxIns prepares utxos for fund tx by calculating fees
func (b *Builder) PrepareFundTxIns() error {
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
func (b *Builder) witsigForFundTxIn(tx *wire.MsgTx) ([]byte, error) {
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

	pub := b.dlc.pubs[b.party]

	return b.wallet.WitnessSignature(tx, fundTxInAt, famt, fc, pub)
}
