package integration

import (
	"testing"

	"github.com/btcsuite/btcutil"
	"github.com/dgarage/dlc/internal/oracle"
	"github.com/stretchr/testify/assert"
)

func contratorHasBalance(t *testing.T, c *Contractor, balance btcutil.Amount) {
	addr, err := c.Wallet.NewAddress()
	assert.NoError(t, err)

	err = Faucet(addr, balance)
	assert.NoError(t, err)
}

func contractorGetCommitmentsFromOracle(t *testing.T, c *Contractor, o *oracle.Oracle) {
	// fixing time of the contract
	fixingTime := c.DLCBuilder.DLC().Conds.FixingTime

	// oracle provides pubkey set for the given time
	pubkeySet, err := o.PubkeySet(fixingTime)
	assert.NoError(t, err)

	// contractor sets and prepare commitents on each deal
	c.DLCBuilder.SetOraclePubkeySet(&pubkeySet)
}

// A contractor sends pubkey and fund txins to the counterparty
func contractorOfferCounterparty(t *testing.T, c1, c2 *Contractor) {
	// first party prepare pubkey and fund txins/txouts
	c1.DLCBuilder.PreparePubkey()
	err := c1.DLCBuilder.PrepareFundTxIns()
	assert.NoError(t, err)

	// send prepared data to second party
	dlc1 := *c1.DLCBuilder.DLC()

	// second party accepts it
	c2.DLCBuilder.CopyReqsFromCounterparty(&dlc1)
}
