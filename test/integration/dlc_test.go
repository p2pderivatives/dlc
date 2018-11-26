package integration

import (
	"math"
	"math/rand"
	"testing"
	"time"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/dgarage/dlc/internal/dlc"
	"github.com/dgarage/dlc/internal/oracle"
	"github.com/dgarage/dlc/internal/rpc"
	"github.com/dgarage/dlc/internal/utils"
	"github.com/stretchr/testify/assert"
)

// Senario: There's a `lottery` oracle who publishes a random digit number everyday,
// and 2 parties bet on tomorrow's numbers randomly
func TestContractorMakeAndExecuteDLC(t *testing.T) {
	// Given an oracle "Olivia"
	nDigit := 2
	olivia, _ := newOracle("Olivia", nDigit)

	// And next announcement time
	fixingTime := nextLotteryAnnouncement()

	// And a contractor "Alice"
	alice, _ := newContractor("Alice")
	balanceA := btcutil.Amount(2 * btcutil.SatoshiPerBitcoin)
	contratorHasBalance(t, alice, balanceA)

	// And a contractor "Bob"
	bob, _ := newContractor("Bob")
	balanceB := btcutil.Amount(2 * btcutil.SatoshiPerBitcoin)
	contratorHasBalance(t, bob, balanceB)

	// -- Making DLC --

	// When Alice and Bob bet on tomorrow's lottery
	contractorsBetOnLottery(
		t, alice, bob, nDigit, fixingTime, blockHeightAfterDays(2))

	// And Alice offers a DLC to Bob
	contractorGetCommitmentsFromOracle(t, alice, olivia)
	contractorOfferCounterparty(t, alice, bob)

	// And Bob accepts the offer
	contractorGetCommitmentsFromOracle(t, bob, olivia)
	contractorAcceptOffer(t, bob, alice)

	// And Alice signs all txs and send the signs to Bob
	contractorSignAllTxs(t, alice, bob)

	// And Bob sends fund tx to the network
	contractorSendFundTx(t, bob)

	// Then Alice and Bob should have remaining balance after funding
	contractorShouldHaveBalanceAfterFunding(t, alice, balanceA)
	contractorShouldHaveBalanceAfterFunding(t, bob, balanceB)

	// -- Executing Contract --

	// When Olivia fixes a number
	oracleFixLottery(t, olivia, nDigit, fixingTime)

	// And Alice and Bob fixe a deal using Olivia's messages and sign
	contractorFixLotteryDeal(t, alice, olivia, nDigit)
	contractorFixLotteryDeal(t, bob, olivia, nDigit)

	// And Alice sends CETx and closing tx
	contractorExecuteContract(t, alice)

	// Alice and Bob should receive funds accoding to the fixed deal
	contractorShouldReceiveFundsByFixedDeal(t, alice, balanceA)
	contractorShouldReceiveFundsByFixedDeal(t, bob, balanceB)
}

// Senario: either party refunds after contract period
func TestContractorRefundDLC(t *testing.T) {
	// Given an oracle "Olivia"
	nDigit := 2
	olivia, _ := newOracle("Olivia", nDigit)

	// And next announcement time
	fixingTime := nextLotteryAnnouncement()

	// And a contractor "Alice"
	alice, _ := newContractor("Alice")
	balanceA := btcutil.Amount(2 * btcutil.SatoshiPerBitcoin)
	contratorHasBalance(t, alice, balanceA)

	// And a contractor "Bob"
	bob, _ := newContractor("Bob")
	balanceB := btcutil.Amount(2 * btcutil.SatoshiPerBitcoin)
	contratorHasBalance(t, bob, balanceB)

	// -- Making DLC --

	// When Alice and Bob bet on tomorrow's lottery
	refundUnlockAt := blockHeightAfterDays(2) // block height of 2 days later
	contractorsBetOnLottery(t, alice, bob, nDigit, fixingTime, refundUnlockAt)

	// And Alice offers a DLC to Bob
	contractorGetCommitmentsFromOracle(t, alice, olivia)
	contractorOfferCounterparty(t, alice, bob)

	// And Bob accepts the offer
	contractorGetCommitmentsFromOracle(t, bob, olivia)
	contractorAcceptOffer(t, bob, alice)

	// And Alice signs all txs and send the signs to Bob
	contractorSignAllTxs(t, alice, bob)

	// And Bob sends fund tx to the network
	contractorSendFundTx(t, bob)

	// Then Alice and Bob should have remaining balance after funding
	contractorShouldHaveBalanceAfterFunding(t, alice, balanceA)
	contractorShouldHaveBalanceAfterFunding(t, bob, balanceB)

	// -- Refunding --

	// When Olivia fixes an invalid messages
	oracleFixInvalidLottery(t, olivia, nDigit, fixingTime)

	// Then Alice and Bob cannot fix deal with invalid messages
	contractorCannotFixLotteryDeal(t, alice, olivia, nDigit)
	contractorCannotFixLotteryDeal(t, bob, olivia, nDigit)

	// And Alice and Bob cannot refund before the unlock time
	contractorCannotRefund(t, alice)
	contractorCannotRefund(t, bob)

	// Given after the unlock time
	waitUntil(t, refundUnlockAt)

	// When alice refund
	contractorRefund(t, alice)

	// Then Alice and Bob should receive funds accoding to the fixed deal
	contractorShouldReceiveRefund(t, alice, balanceA)
	contractorShouldReceiveRefund(t, bob, balanceB)
}

// nextLotteryAnnouncement returns time of tomorrow noon
func nextLotteryAnnouncement() time.Time {
	tomorrow := time.Now().AddDate(0, 0, 1)
	year, month, day := tomorrow.Date()
	return time.Date(year, month, day, 12, 0, 0, 0, tomorrow.Location())
}

func blockHeightAfterDays(days int) uint32 {
	unlockAt := time.Now().AddDate(0, 0, days)
	d := time.Until(unlockAt)
	timePerBlock := chaincfg.RegressionNetParams.TargetTimePerBlock
	blocks := uint32(d / timePerBlock)
	curHeight, _ := rpc.GetBlockCount()
	return uint32(curHeight) + blocks
}

// The both contractors agree on random n-digit number lottery
// by creating all patterns randomly
func contractorsBetOnLottery(
	t *testing.T, c1, c2 *Contractor, nDigit int, fixingTime time.Time, refundUnlockAt uint32) {

	var onebtc btcutil.Amount = 1 * btcutil.SatoshiPerBitcoin
	famta, famtb := onebtc, onebtc
	deals := randomDealsForAllDigitPatterns(nDigit, int(famta+famtb))
	conds, err := dlc.NewConditions(
		fixingTime, famta, famtb, 1, 1, refundUnlockAt, deals)
	assert.NoError(t, err)

	c1.createDLCBuilder(conds, dlc.FirstParty)
	c2.createDLCBuilder(conds, dlc.SecondParty)
}

func randomDealsForAllDigitPatterns(nDigit, famt int) []*dlc.Deal {
	var deals []*dlc.Deal
	for d := 0; d < int(math.Pow10(nDigit)); d++ {
		// randomly decide deals
		damta, damtb := randomDealAmts(famt)
		deal := dlc.NewDeal(damta, damtb, nDigitToBytes(d, nDigit))
		deals = append(deals, deal)
	}
	return deals
}

func randomDealAmts(famt int) (btcutil.Amount, btcutil.Amount) {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	a := r.Intn(famt + 1)
	b := famt - a
	return utils.ItoAmt(a), utils.ItoAmt(b)
}

// convert a n-digit number to byte
// e.g. 123 -> [][]byte{{1}, {2}, {3}}
func nDigitToBytes(d int, n int) [][]byte {
	b := make([][]byte, n)
	for i := 0; i < n; i++ {
		b[i] = []byte{byte(d % 10)}
		d = d / 10
	}
	return b
}

func oracleFixLottery(
	t *testing.T, o *oracle.Oracle, n int, ftime time.Time) {
	msgs := nDigitToBytes(randomNumber(n), n)

	err := o.FixMsgs(ftime, msgs)
	assert.NoError(t, err)
}

func randomNumber(nDigit int) (d int) {
	for i := 0; i < nDigit; i++ {
		d = d + randomDigit()*int(math.Pow10(i))
	}
	return
}

func randomDigit() int {
	s := rand.NewSource(time.Now().UnixNano())
	return rand.New(s).Intn(10)
}

func contractorFixLotteryDeal(
	t *testing.T, c *Contractor, o *oracle.Oracle, n int) {

	// use all messages that oracle publishes
	idxs := []int{}
	for i := 0; i < n; i++ {
		idxs = append(idxs, i)
	}

	contractorFixDeal(t, c, o, idxs)
}

func contractorCannotFixLotteryDeal(
	t *testing.T, c *Contractor, o *oracle.Oracle, n int) {

	// use all messages that oracle publishes
	idxs := []int{}
	for i := 0; i < n; i++ {
		idxs = append(idxs, i)
	}

	contractorCannotFixDeal(t, c, o, idxs)
}

func oracleFixInvalidLottery(
	t *testing.T, o *oracle.Oracle, n int, ftime time.Time) {

	msgs := randomStringMsgs(n)

	err := o.FixMsgs(ftime, msgs)
	assert.NoError(t, err)
}

func randomStringMsgs(n int) [][]byte {
	msgs := make([][]byte, n)
	for i := 0; i < n; i++ {
		msgs[i] = []byte(string(randomChar()))
	}
	return msgs
}

func randomChar() rune {
	var runes = []rune("abcdefghijklmnopqrstuvwxyz")
	s := rand.NewSource(time.Now().UnixNano())
	return runes[rand.New(s).Intn(len(runes))]
}
