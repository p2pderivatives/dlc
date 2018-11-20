package integration

import (
	"log"
	"math"
	"math/rand"
	"os"
	"runtime/pprof"
	"testing"
	"time"

	"github.com/btcsuite/btcutil"
	"github.com/dgarage/dlc/internal/dlc"
	"github.com/dgarage/dlc/internal/oracle"
	"github.com/dgarage/dlc/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestContractorMakeDLC(t *testing.T) {
	cpuprofile := "mycpu.prof"
	f, err := os.Create(cpuprofile)
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	// Given an oracle "Olivia" who publishes random a 5-digit number everday like lottery
	nDigit := 2
	olivia, _ := newOracle("Olivia", nDigit)

	// And next announcement time is noon tomorrow
	fixingTime := nextLotteryAnnouncement()

	// And a contractor "Alice"
	alice, _ := newContractor("Alice")

	// And a contractor "Bob"
	bob, _ := newContractor("Bob")

	// And Alice and Bob bet on all cases
	contractorsBetOnAllDigitPatters(t, alice, bob, nDigit, fixingTime)

	// When Alice offers a DLC to Bob
	contractorGetCommitmentsFromOracle(t, alice, olivia)
	contractorOfferCounterparty(t, alice, bob)
}

func nextLotteryAnnouncement() time.Time {
	tomorrow := time.Now().AddDate(0, 0, 1)
	year, month, day := tomorrow.Date()
	return time.Date(year, month, day, 12, 0, 0, 0, tomorrow.Location())
}

// The both contractors agree on conditions of DLC
// In this case, random 5-digit numbers and random deals
func contractorsBetOnAllDigitPatters(
	t *testing.T, c1, c2 *Contractor, nDigit int, fixingTime time.Time) {

	var onebtc btcutil.Amount = 1 * btcutil.SatoshiPerBitcoin
	famta, famtb := onebtc, onebtc
	deals := randomDealsForAllDigitPatterns(nDigit, int(famta+famtb))
	conds, err := dlc.NewConditions(fixingTime, famta, famtb, 1, 1, 1, deals)
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
	s := rand.NewSource(time.Now().Unix())
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
