package integration

import (
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/p2pderivatives/dlc/internal/oracle"
	"github.com/p2pderivatives/dlc/pkg/dlc"
	"github.com/stretchr/testify/assert"
)

func TestOracleCommitAndSign(t *testing.T) {
	// Given an oracle "Olivia" who provides weather information
	//   weather info contains "weather" "temperature" "windspeed"
	olivia, _ := newOracle("Olivia", 3)

	// And a contractor "Alice"
	alice, _ := newContractor("Alice")

	// And Alice bet on "weather" and "temprature" at a future time
	fixingTime := contractorBetOnWeatherAndTemperature(t, alice)

	// Alice asks Olivia to fix weather info at the fixing time
	contractorAsksOracleToCommit(t, alice, olivia)

	// Olivia fixes weather info
	fixedWeather := oracleFixesWeather(t, olivia, fixingTime)

	// When Alice fixes a deal using Olivia's sign and singed weather info
	contractorFixesWeatherDeal(t, alice, olivia)

	// Then The fixed deal should be a subset of the fixed weather info
	shouldFixedDealSameWithFixedWeather(t, alice, fixedWeather)
}

func newWeather(weather string, temp int, windSpeed int) [][]byte {
	return [][]byte{
		[]byte(weather),
		[]byte(strconv.Itoa(temp)),
		[]byte(strconv.Itoa(windSpeed)),
	}
}

func contractorBetOnWeatherAndTemperature(t *testing.T, c *Contractor) time.Time {
	deal1 := dlc.NewDeal(2, 0, newWeather("fine", 20, 0)[:2])
	deal2 := dlc.NewDeal(1, 1, newWeather("fine", 10, 0)[:2])
	deal3 := dlc.NewDeal(1, 1, newWeather("rain", 20, 0)[:2])
	deal4 := dlc.NewDeal(0, 2, newWeather("rain", 10, 0)[:2])
	deals := []*dlc.Deal{deal1, deal2, deal3, deal4}
	fixingTime := time.Now().AddDate(0, 0, 1)
	net := &chaincfg.RegressionNetParams
	conds, err := dlc.NewConditions(net, fixingTime, 1, 1, 1, 1, 1, deals)
	assert.NoError(t, err)
	c.createDLCBuilder(conds, dlc.FirstParty)
	return fixingTime
}

func contractorAsksOracleToCommit(
	t *testing.T, c *Contractor, o *oracle.Oracle) {
	ftime := c.DLCBuilder.Contract.Conds.FixingTime

	pubkeySet, err := o.PubkeySet(ftime)
	assert.NoError(t, err)

	err = c.DLCBuilder.SetOraclePubkeySet(&pubkeySet)
	assert.NoError(t, err)
}

func oracleFixesWeather(
	t *testing.T, o *oracle.Oracle, ftime time.Time) [][]byte {
	msgs := [][][]byte{
		newWeather("fine", 20, 0),
		newWeather("fine", 20, 5),
		newWeather("fine", 10, 0),
		newWeather("fine", 10, 5),
		newWeather("rain", 20, 0),
		newWeather("rain", 20, 5),
		newWeather("rain", 10, 0),
		newWeather("rain", 10, 5),
	}

	fixingMsg := randomMsg(msgs)

	err := o.FixMsgs(ftime, fixingMsg)
	assert.NoError(t, err)

	return fixingMsg
}

func randomMsg(msgs [][][]byte) [][]byte {
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)
	idx := r.Intn(len(msgs))
	return msgs[idx]
}

func contractorFixesWeatherDeal(t *testing.T, c *Contractor, o *oracle.Oracle) {
	idxs := []int{0, 1} // use only weather and temperature
	contractorFixDeal(t, c, o, idxs)
}

func shouldFixedDealSameWithFixedWeather(t *testing.T, c *Contractor, fixedWeather [][]byte) {
	_, fixedDeal, err := c.DLCBuilder.Contract.FixedDeal()
	assert.NoError(t, err)
	assert.Equal(t, fixedWeather[:2], fixedDeal.Msgs)
}
