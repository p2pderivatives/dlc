package dlc

import (
	"fmt"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/dgarage/dlc/internal/wallet"
)

// Tx sizes for fee estimation
const fundTxBaseSize = int64(55)
const fundTxInSize = int64(149)
const fundTxOutSize = int64(31)

// Builder builds DLC by interacting with wallet
type Builder struct {
	party         Contractor
	wallet        *wallet.Wallet
	dlc           *DLC
	feeCalculator FeeCalculator
}

// FeeCalculator calculates fee in sathoshi based on bytes
type FeeCalculator func(int64) int64

// NewBuilder creates a new Builder for a contractor
func NewBuilder(party Contractor, wallet *wallet.Wallet, feeCalc FeeCalculator) *Builder {
	return &Builder{party: party, wallet: wallet, feeCalculator: feeCalc}
}

// CreateDraft creates a DLC draft based on given conditions
func (b Builder) CreateDraft(fundAmts FundAmounts) {
	b.dlc = &DLC{fundAmts: fundAmts}
}

// PrepareFundTx prepares utxos for fund tx by calculating fees
func (b Builder) PrepareFundTx() error {
	famt := b.dlc.fundAmts.Amount(b.party)
	feeBase := b.feeCalculator(fundTxBaseSize)

	utxos, change, err := b.selectUtxos(famt + feeBase)
	if err != nil {
		return err
	}

	txins, err := utxosToTxIns(utxos)
	if err != nil {
		return err
	}
	b.dlc.SetFundTxIns(b.party, txins)

	if change > 0 {
		pkScript, err := b.wallet.NewWitnessPubkeyScript()
		if err != nil {
			return err
		}
		txout := wire.NewTxOut(change, pkScript)
		b.dlc.SetFundTxOuts(b.party, []*wire.TxOut{txout})
	}

	return nil
}

// selectUtxos selects utxos for requested amount
func (b Builder) selectUtxos(amt int64) (utxos []wallet.Utxo, change int64, err error) {
	var allUtxos []wallet.Utxo
	allUtxos, err = b.wallet.ListUnspent()
	if err != nil {
		return
	}

	var total int64
	var fee int64
	var utxoAmt btcutil.Amount
	for _, utxo := range allUtxos {
		utxoAmt, err = btcutil.NewAmount(utxo.Amount)
		if err != nil {
			return
		}
		total += int64(utxoAmt)
		fee += b.feeCalculator(fundTxInSize)
		utxos = append(utxos, utxo)
		if amt+fee == total {
			return
		} else if amt+fee < total {
			change = total - (amt + fee)
			fee += b.feeCalculator(fundTxOutSize)
			if amt+fee <= total {
				return
			}
		}
	}

	err = fmt.Errorf("Not enough utxos")
	return
}

func utxosToTxIns(utxos []wallet.Utxo) ([]*wire.TxIn, error) {
	var txins []*wire.TxIn
	for _, utxo := range utxos {
		txid, err := chainhash.NewHashFromStr(utxo.TxID)
		if err != nil {
			return txins, err
		}
		op := wire.NewOutPoint(txid, utxo.Vout)
		txins = append(txins, wire.NewTxIn(op, nil, nil))
	}
	return txins, nil
}
