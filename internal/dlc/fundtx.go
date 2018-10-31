package dlc

import "github.com/btcsuite/btcd/wire"

// FundAmounts contains fund amounts in satoshi
type FundAmounts struct {
	amts map[Contractor]int64
}

// NewFundAmounts creates FundAmounts
func NewFundAmounts(amt1, amt2 int64) FundAmounts {
	amts := make(map[Contractor]int64)
	amts[FirstParty] = amt1
	amts[SecondParty] = amt2
	return FundAmounts{amts: amts}
}

// Amount returns fund amount for a party
func (amts *FundAmounts) Amount(party Contractor) int64 {
	return amts.amts[party]
}

// FundTxRequirements contains txins and txouts for fund tx
type FundTxRequirements struct {
	txIns  map[Contractor]int64
	txOuts map[Contractor]int64
}

// SetFundTxIns sets txins of a party
func (dlc *DLC) SetFundTxIns(party Contractor, txins []*wire.TxIn) {
	dlc.fundTxReqs.txIns[party] = txins
}

// SetFundTxOuts sets txouts of a party
func (dlc *DLC) SetFundTxOuts(party Contractor, txouts []*wire.TxOut) {
	dlc.fundTxReqs.txOuts[party] = txouts
}
