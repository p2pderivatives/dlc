package dlc

import ()

// DLC contains all information required for DLC contract
// including FundTx, SettlementTx, RefundTx
type DLC struct {
	fundAmts   FundAmounts
	fundTxReqs *FundTxRequirements
}
