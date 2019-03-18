package dlc

import (
	"github.com/btcsuite/btcutil"
)

// Tx sizes for fee estimation
const fundTxBaseSize = int64(55)
const fundTxInSize = int64(149)
const fundTxOutSize = int64(31)
const cetxSize = int64(345) // context execution tx size
const closingTxSize = int64(238)

func (d *DLC) fundTxFeeBase() btcutil.Amount {
	return d.Conds.FundFeerate.MulF64(float64(fundTxBaseSize))
}

func (d *DLC) fundTxFeeTxIns(n int) btcutil.Amount {
	return d.Conds.FundFeerate.MulF64(float64(fundTxInSize * int64(n)))
}

func (d *DLC) fundTxFeePerTxIn() btcutil.Amount {
	return d.Conds.FundFeerate.MulF64(float64(fundTxInSize))
}

func (d *DLC) fundTxFeePerTxOut() btcutil.Amount {
	return d.Conds.FundFeerate.MulF64(float64(fundTxOutSize))
}

func (d *DLC) fundTxFeeByParty(p Contractor) btcutil.Amount {
	feeBase := d.fundTxFeeBase() / 2
	feeIns := d.fundTxFeeTxIns(len(d.Utxos[p]))
	feeOut := btcutil.Amount(0)
	if d.ChangeAddrs[p] != nil {
		feeOut = d.fundTxFeePerTxOut()
	}
	return feeBase + feeIns + feeOut
}

func (d *DLC) redeemTxFee(size int64) btcutil.Amount {
	return d.Conds.RedeemFeerate.MulF64(float64(size))
}

func (d *DLC) feeByParty(p Contractor) btcutil.Amount {
	fundtxFee := d.fundTxFeeByParty(p)
	cetxFee := d.redeemTxFee(cetxSize) / 2
	closingtxFee := d.redeemTxFee(closingTxSize) / 2
	return fundtxFee + cetxFee + closingtxFee
}
