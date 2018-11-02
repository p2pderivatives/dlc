package dlc

import (
	"fmt"

	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/dgarage/dlc/internal/wallet"
)

// SetFundAmounts sets fund amounts to DLC
func (b *Builder) SetFundAmounts(amt1, amt2 int64) {
	b.dlc.fundAmts[FirstParty] = amt1
	b.dlc.fundAmts[SecondParty] = amt2
}

func (b *Builder) fundAmountByParty(party Contractor) (int64, error) {
	famt, ok := b.dlc.fundAmts[party]
	if !ok {
		err := fmt.Errorf("fund amount isn't set yet")
		return 0, err
	}
	return famt, nil
}

// FundTxRequirements contains txins and txouts for fund tx
type FundTxRequirements struct {
	txIns  map[Contractor][]*wire.TxIn
	txOuts map[Contractor][]*wire.TxOut
}

func (b *Builder) setFundTxIns(party Contractor, txins []*wire.TxIn) {
	b.dlc.fundTxReqs.txIns[party] = txins
}

func (b *Builder) setFundTxOuts(party Contractor, txouts []*wire.TxOut) {
	b.dlc.fundTxReqs.txOuts[party] = txouts
}

// Tx sizes for fee estimation
const fundTxBaseSize = int64(55)
const fundTxInSize = int64(149)
const fundTxOutSize = int64(31)

// PrepareFundTx prepares utxos for fund tx by calculating fees
func (b *Builder) PrepareFundTx() error {
	famt, err := b.fundAmountByParty(b.party)
	if err != nil {
		return err
	}
	feeBase := b.feeCalc(fundTxBaseSize)

	utxos, change, err := b.selectUtxos(famt + feeBase)
	if err != nil {
		return err
	}

	txins, err := wallet.UtxosToTxIns(utxos)
	if err != nil {
		return err
	}
	b.setFundTxIns(b.party, txins)

	if change > 0 {
		pkScript, err := b.wallet.NewWitnessPubkeyScript()
		if err != nil {
			return err
		}
		txout := wire.NewTxOut(change, pkScript)
		b.setFundTxOuts(b.party, []*wire.TxOut{txout})
	}

	return nil
}

// selectUtxos selects utxos for requested amount
// TODO: move utxo selection logic to wallet package by removing dependencies on tx sizes
func (b *Builder) selectUtxos(amt int64) (utxos []wallet.Utxo, change int64, err error) {
	var utxosAll []wallet.Utxo
	utxosAll, err = b.wallet.ListUnspent()
	if err != nil {
		return
	}

	var total int64
	var fee int64
	var utxoAmt btcutil.Amount
	for _, utxo := range utxosAll {
		utxoAmt, err = btcutil.NewAmount(utxo.Amount)
		if err != nil {
			return
		}
		total += int64(utxoAmt)
		fee += b.feeCalc(fundTxInSize)
		utxos = append(utxos, utxo)
		if amt+fee == total {
			return
		} else if amt+fee < total {
			change = total - (amt + fee)
			fee += b.feeCalc(fundTxOutSize)
			if amt+fee <= total {
				return
			}
		}
	}

	err = fmt.Errorf("Not enough utxos")
	return
}

const fundTxVersion = 2

// FundTx constracts fund tx using prepared fund tx requirements
func (dlc *DLC) FundTx() *wire.MsgTx {
	tx := wire.NewMsgTx(fundTxVersion)
	return tx
}
