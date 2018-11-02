package dlc

import (
	"github.com/dgarage/dlc/internal/wallet"
)

// DLC contains all information required for DLC contract
// including FundTx, SettlementTx, RefundTx
type DLC struct {
	fundAmts   map[Contractor]int64
	fundTxReqs *FundTxRequirements
}

func newDLC() *DLC {
	return &DLC{
		fundAmts:   make(map[Contractor]int64),
		fundTxReqs: newFundTxRequirements(),
	}
}

// Contractor represents a contractor type
type Contractor int

const (
	// FirstParty is a contractor who creates offer
	FirstParty Contractor = 0
	// SecondParty is a contractor who accepts offer
	SecondParty Contractor = 1
)

// Builder builds DLC by interacting with wallet
type Builder struct {
	party   Contractor
	wallet  wallet.Wallet
	dlc     *DLC
	feeCalc FeeCalculator
}

// FeeCalculator calculates fee in sathoshi based on bytes
type FeeCalculator func(int64) int64

// NewBuilder creates a new Builder for a contractor
func NewBuilder(
	party Contractor, w wallet.Wallet, feeCalc FeeCalculator,
) *Builder {
	return &Builder{
		dlc:     newDLC(),
		party:   party,
		wallet:  w,
		feeCalc: feeCalc,
	}
}

// DLC returns the DLC constructed by builder
func (b *Builder) DLC() *DLC {
	return b.dlc
}
