package dlc

import (
	"fmt"

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

func newFundTxRequirements() *FundTxRequirements {
	return &FundTxRequirements{
		txIns: make(map[Contractor][]*wire.TxIn),
		txOut: make(map[Contractor]*wire.TxOut),
	}
}

const fundTxVersion = 2

// FundTx constructs fund tx using prepared fund tx requirements
func (d *DLC) FundTx() *wire.MsgTx {
	tx := wire.NewMsgTx(fundTxVersion)

	// TODO: add txout script for the txin of settlement tx

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

	return tx
}

// SetFundAmounts sets fund amounts to DLC
func (b *Builder) SetFundAmounts(amt1, amt2 btcutil.Amount) {
	b.dlc.fundAmts[FirstParty] = amt1
	b.dlc.fundAmts[SecondParty] = amt2
}

// SetFundFeerate sets feerate (satoshi/byte) for fund tx fee calculation
func (b *Builder) SetFundFeerate(feerate btcutil.Amount) {
	b.dlc.fundFeerate = feerate
}

// Tx sizes for fee estimation
const fundTxBaseSize = int64(55)
const fundTxInSize = int64(149)
const fundTxOutSize = int64(31)

func (d *DLC) fundTxFeeBase() btcutil.Amount {
	return d.fundFeerate.MulF64(float64(fundTxBaseSize))
}

func (d *DLC) fundTxFeePerTxIn() btcutil.Amount {
	return d.fundFeerate.MulF64(float64(fundTxInSize))
}

func (d *DLC) fundTxFeePerTxOut() btcutil.Amount {
	return d.fundFeerate.MulF64(float64(fundTxOutSize))
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
		b.dlc.fundTxReqs.txOut[b.party] = txout
	}

	return nil
}
