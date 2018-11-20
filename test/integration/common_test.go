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

// A contractor sends pubkey, fund txins and
// signs of context execution txs and refund tx
func contractorAcceptOffer(t *testing.T, c1, c2 *Contractor) {
	// Second party prepares pubkey and fund txins/txouts
	c1.DLCBuilder.PreparePubkey()
	err := c1.DLCBuilder.PrepareFundTxIns()
	assert.NoError(t, err)

	//  signs txs
	ceSigns, rfSign := conractorSignCETxsAndRefundTx(t, c1)

	// Sends pubkey and fund txins and sign to the counterparty
	dlc1 := *c1.DLCBuilder.DLC()
	c2.DLCBuilder.CopyReqsFromCounterparty(&dlc1)
	contractorSendTxSigns(t, c2, ceSigns, rfSign)
}

// A contractor sends witnesses of fund txins to allow the counterparty to send fund tx
func contractorSendFundTxInWitnesses(t *testing.T, c1, c2 *Contractor) {
	// signs
	ceSigns, rfSign := conractorSignCETxsAndRefundTx(t, c1)
	contractorSendTxSigns(t, c2, ceSigns, rfSign)

	// TODO: send fund txin witness to the counterparty
}

// A contractor signs CETxs and Refund tx and sends them to the counterparty
func conractorSignCETxsAndRefundTx(t *testing.T, c *Contractor) ([][]byte, []byte) {
	// unlocks to sign txs
	c.unlockWallet()

	// context execution txs signs
	ceSigns, err := c.DLCBuilder.SignContractExecutionTxs()
	assert.NoError(t, err)

	// create refund tx sign
	rfSign, err := c.DLCBuilder.SignRefundTx()
	assert.NoError(t, err)

	return ceSigns, rfSign
}

func contractorSendTxSigns(t *testing.T, c *Contractor, ceSigns [][]byte, rfSign []byte) {
	err := c.DLCBuilder.AcceptCETxSigns(ceSigns)
	assert.NoError(t, err)
	err = c.DLCBuilder.AcceptRefundTxSign(rfSign)
	assert.NoError(t, err)
}

func contractorSendFundTx(t *testing.T, c *Contractor) {
	err := c.DLCBuilder.SendFundTx()
	assert.NoError(t, err)
}
