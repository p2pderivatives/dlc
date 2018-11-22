package integration

import (
	"testing"

	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/dgarage/dlc/internal/oracle"
	"github.com/dgarage/dlc/internal/rpc"
	"github.com/stretchr/testify/assert"
)

func contratorHasBalance(t *testing.T, c *Contractor, balance btcutil.Amount) {
	addr, err := c.Wallet.NewAddress()
	assert.NoError(t, err)

	err = rpc.Faucet(addr, balance)
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

// A contractor sends pubkey, fund txins and
// signs of context execution txs and refund tx
func contractorAcceptOffer(t *testing.T, c1, c2 *Contractor) {
	// Second party prepares pubkey and fund txins/txouts
	c1.DLCBuilder.PreparePubkey()
	err := c1.DLCBuilder.PrepareFundTxIns()
	assert.NoError(t, err)

	// signs CE txs and refund tx
	ceSigns := conractorSignCETxs(t, c1)
	rfSign := conractorSignRefundTx(t, c1)

	// Sends pubkey and fund txins and sign to the counterparty
	dlc1 := *c1.DLCBuilder.DLC()
	c2.DLCBuilder.CopyReqsFromCounterparty(&dlc1)

	// send signs
	err = c2.DLCBuilder.AcceptCETxSigns(ceSigns)
	assert.NoError(t, err)
	err = c2.DLCBuilder.AcceptRefundTxSign(rfSign)
	assert.NoError(t, err)
}

// A contractor sends signs of all transactions (fund tx, CE txs, refund tx)
func contractorSignAllTxs(t *testing.T, c1, c2 *Contractor) {
	// signs all txs
	ceSigns := conractorSignCETxs(t, c1)
	rfSign := conractorSignRefundTx(t, c1)
	fundWits := contractorSignFundTx(t, c1)

	// send all signs and witnesses
	err := c2.DLCBuilder.AcceptCETxSigns(ceSigns)
	assert.NoError(t, err)
	err = c2.DLCBuilder.AcceptRefundTxSign(rfSign)
	assert.NoError(t, err)
	c2.DLCBuilder.AcceptFundWitnesses(fundWits)
}

// A contractor signs CETxs
func conractorSignCETxs(t *testing.T, c *Contractor) [][]byte {
	// unlocks to sign txs
	c.unlockWallet()

	// context execution txs signs
	ceSigns, err := c.DLCBuilder.SignContractExecutionTxs()
	assert.NoError(t, err)

	return ceSigns
}

// A contractor signs Refund tx
func conractorSignRefundTx(t *testing.T, c *Contractor) []byte {
	// unlocks to sign txs
	c.unlockWallet()

	// create refund tx sign
	rfSign, err := c.DLCBuilder.SignRefundTx()
	assert.NoError(t, err)

	return rfSign
}

func contractorSignFundTx(t *testing.T, c *Contractor) []wire.TxWitness {
	// unlocks to sign txs
	c.unlockWallet()

	// create fund tx witnesses
	wits, err := c.DLCBuilder.SignFundTx()
	assert.NoError(t, err)

	return wits
}

func contractorSendFundTx(t *testing.T, c *Contractor) {
	_, err := c.DLCBuilder.SignFundTx()
	assert.NoError(t, err)
	err = c.DLCBuilder.SendFundTx()
	assert.NoError(t, err)
}
