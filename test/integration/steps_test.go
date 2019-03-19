package integration

import (
	"testing"

	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/p2pderivatives/dlc/internal/oracle"
	"github.com/p2pderivatives/dlc/internal/rpc"
	"github.com/stretchr/testify/assert"
)

var feeByParty = btcutil.Amount(float64(415+345+238) / 2)

func contratorHasBalance(t *testing.T, c *Contractor, balance btcutil.Amount) {
	addr, err := c.Wallet.NewAddress()
	assert.NoError(t, err)

	err = rpc.Faucet(addr, balance)
	assert.NoError(t, err)
}

func contractorGetCommitmentsFromOracle(t *testing.T, c *Contractor, o *oracle.Oracle) {
	// fixing time of the contract
	fixingTime := c.DLCBuilder.Contract.Conds.FixingTime

	// oracle provides pubkey set for the given time
	pubkeySet, err := o.PubkeySet(fixingTime)
	assert.NoError(t, err)

	// contractor sets and prepare commitents on each deal
	nRpoints := len(pubkeySet.CommittedRpoints)
	idxs := []int{}
	for idx := 0; idx < nRpoints; idx++ {
		idxs = append(idxs, idx)
	}
	err = c.DLCBuilder.SetOraclePubkeySet(&pubkeySet, idxs)
	assert.NoError(t, err)
}

// A contractor sends pubkey and fund txins to the counterparty
func contractorOfferCounterparty(t *testing.T, c1, c2 *Contractor) {
	// first party prepare pubkey and fund txins/txouts
	c1.DLCBuilder.PreparePubkey()
	err := c1.DLCBuilder.PrepareFundTx()
	assert.NoError(t, err)

	// send pubkey adn utxo to the counterparty
	p1, err := c1.DLCBuilder.PublicKey()
	assert.NoError(t, err)
	u1 := c1.DLCBuilder.Utxos()
	addr1 := c1.DLCBuilder.Address()
	caddr1 := c1.DLCBuilder.ChangeAddress()

	// the counterparty party accepts them
	err = c2.DLCBuilder.AcceptPubkey(p1)
	assert.NoError(t, err)
	err = c2.DLCBuilder.AcceptUtxos(u1)
	assert.NoError(t, err)
	c2.DLCBuilder.AcceptAdderss(addr1)
	c2.DLCBuilder.AcceptChangeAdderss(caddr1)
}

// A contractor sends pubkey, fund txins and
// signatures of context execution txs and refund tx
func contractorAcceptOffer(t *testing.T, c1, c2 *Contractor) {
	// Second party prepares pubkey and fund txins/txouts
	c1.DLCBuilder.PreparePubkey()
	err := c1.DLCBuilder.PrepareFundTx()
	assert.NoError(t, err)

	// signatures CE txs and refund tx
	ceSigs := conractorSignCETxs(t, c1)
	rfSig := conractorSignRefundTx(t, c1)

	// sends pubkey and fund txins and sign to the counterparty
	p1, err := c1.DLCBuilder.PublicKey()
	assert.NoError(t, err)
	u1 := c1.DLCBuilder.Utxos()
	addr1 := c1.DLCBuilder.Address()
	caddr1 := c1.DLCBuilder.ChangeAddress()

	// the counterparty accepts them
	err = c2.DLCBuilder.AcceptPubkey(p1)
	assert.NoError(t, err)
	err = c2.DLCBuilder.AcceptUtxos(u1)
	assert.NoError(t, err)
	c2.DLCBuilder.AcceptAdderss(addr1)
	c2.DLCBuilder.AcceptChangeAdderss(caddr1)

	// send sigatures
	err = c2.DLCBuilder.AcceptCETxSignatures(ceSigs)
	assert.NoError(t, err)
	err = c2.DLCBuilder.AcceptRefundTxSignature(rfSig)
	assert.NoError(t, err)
}

// A contractor sends signs of all transactions (fund tx, CE txs, refund tx)
func contractorSignAllTxs(t *testing.T, c1, c2 *Contractor) {
	// signs all txs
	ceSigs := conractorSignCETxs(t, c1)
	rfSig := conractorSignRefundTx(t, c1)
	fundWits := contractorSignFundTx(t, c1)

	// send all signs and witnesses
	err := c2.DLCBuilder.AcceptCETxSignatures(ceSigs)
	assert.NoError(t, err)
	err = c2.DLCBuilder.AcceptRefundTxSignature(rfSig)
	assert.NoError(t, err)
	c2.DLCBuilder.AcceptFundWitnesses(fundWits)
}

// A contractor signs CETxs
func conractorSignCETxs(t *testing.T, c *Contractor) [][]byte {
	// unlocks to sign txs
	c.unlockWallet()

	// context execution txs signatures
	ceSigs, err := c.DLCBuilder.SignContractExecutionTxs()
	assert.NoError(t, err)

	return ceSigs
}

// A contractor signs Refund tx
func conractorSignRefundTx(t *testing.T, c *Contractor) []byte {
	// unlocks to sign txs
	c.unlockWallet()

	// create refund tx sign
	rfSig, err := c.DLCBuilder.SignRefundTx()
	assert.NoError(t, err)

	return rfSig
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

	_, err = rpc.Generate(1)
	assert.NoError(t, err)
}

func contractorShouldHaveBalanceAfterFunding(
	t *testing.T, c *Contractor, balanceBefore btcutil.Amount) {
	fundAmt := c.DLCBuilder.FundAmt()
	balance, err := c.balance()
	assert.NoError(t, err)

	// expected_balance = balance_before - fund_amount - fee
	expected := int64(balanceBefore - fundAmt - feeByParty)
	actual := int64(balance)
	assert.InDelta(t, expected, actual, 1)
}

func contractorFixDeal(
	t *testing.T, c *Contractor, o *oracle.Oracle, idxs []int) {
	ftime := c.DLCBuilder.Contract.Conds.FixingTime

	// receive signed message
	sm, err := o.SignMsg(ftime)
	assert.NoError(t, err)

	// fix deal with the fixed msg
	err = c.DLCBuilder.FixDeal(&sm, idxs)
	assert.NoError(t, err)
}

func contractorCannotFixDeal(
	t *testing.T, c *Contractor, o *oracle.Oracle, idxs []int) {
	ftime := c.DLCBuilder.Contract.Conds.FixingTime

	// receive signset
	signedMsg, err := o.SignMsg(ftime)
	assert.NoError(t, err)

	// fail to fix deal with the siged msg
	err = c.DLCBuilder.FixDeal(&signedMsg, idxs)
	assert.Error(t, err)
}

func contractorExecuteContract(t *testing.T, c *Contractor) {
	err := c.DLCBuilder.ExecuteContract()
	assert.NoError(t, err)

	_, err = rpc.Generate(1)
	assert.NoError(t, err)
}

func contractorShouldReceiveFundsByFixedDeal(
	t *testing.T, c *Contractor, balanceBefore btcutil.Amount) {

	fundAmt := c.DLCBuilder.FundAmt()
	dealAmt, err := c.DLCBuilder.FixedDealAmt()
	assert.NoError(t, err)
	balance, err := c.balance()
	assert.NoError(t, err)

	// expected_balance =
	//   balance_before - fund_amount + deal_amount - fee
	expected := int64(balanceBefore - fundAmt + dealAmt - feeByParty)
	actual := int64(balance)
	assert.InDelta(t, expected, actual, 1)
}

func contractorRefund(t *testing.T, c *Contractor) {
	err := c.DLCBuilder.SendRefundTx()
	assert.NoError(t, err)

	_, err = rpc.Generate(1)
	assert.NoError(t, err)
}

func contractorCannotRefund(t *testing.T, c *Contractor) {
	err := c.DLCBuilder.SendRefundTx()
	assert.Error(t, err)
}

func contractorShouldReceiveRefund(
	t *testing.T, c *Contractor, balanceBefore btcutil.Amount) {

	balance, err := c.balance()
	assert.NoError(t, err)

	// expected_balance = balance_before - fee
	expected := int64(balanceBefore - feeByParty)
	actual := int64(balance)

	assert.InDelta(t, expected, actual, 1)
}

func waitUntil(t *testing.T, height uint32) {
	curHeight, err := rpc.GetBlockCount()
	assert.NoError(t, err)

	_, err = rpc.Generate(uint32(int64(height) - curHeight))
	assert.NoError(t, err)
}
