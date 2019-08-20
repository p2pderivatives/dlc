package dlc

import (
	"testing"
	"time"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/stretchr/testify/assert"
)

var testAddress = "bcrt1q8cjx85nnuqd92mq3xnfrqc4xxljhm5sjax55rk"

func TestCondions(t *testing.T) {
	assert := assert.New(t)

	net := &chaincfg.RegressionNetParams
	ftime := time.Now().AddDate(0, 0, 1)
	var famt1, famt2,
		frate, rrate btcutil.Amount = 1, 1, 1, 1
	var lc uint32 = 1
	deals := []*Deal{NewDeal(1, 1, [][]byte{{1}})}

	var err error
	_, err = NewConditions(
		net, ftime, famt1, famt2, frate, rrate, lc, deals, nil)
	assert.NoError(err)

	_, err = NewConditions(
		net, time.Now(), famt1, famt2, frate, rrate, lc, deals, nil)
	assert.Error(err)

	_, err = NewConditions(
		net, ftime, 0, famt2, frate, rrate, lc, deals, nil)
	assert.NoError(err)

	_, err = NewConditions(
		net, ftime, famt1, 0, frate, rrate, lc, deals, nil)
	assert.NoError(err)

	_, err = NewConditions(
		net, ftime, 0, 0, frate, rrate, lc, deals, nil)
	assert.Error(err)

	_, err = NewConditions(
		net, ftime, famt1, famt2, 0, rrate, lc, deals, nil)
	assert.Error(err)

	_, err = NewConditions(
		net, ftime, famt1, famt2, frate, 0, lc, deals, nil)
	assert.Error(err)

	_, err = NewConditions(
		net, ftime, famt1, famt2, frate, rrate, lc, []*Deal{}, nil)
	assert.Error(err)
}

func TestNewPremiumInfo_CorrectParameters_NoError(t *testing.T) {
	// Arrange
	assert := assert.New(t)
	address, err := btcutil.DecodeAddress(testAddress, &chaincfg.RegressionNetParams)
	assert.NoError(err)

	amount := btcutil.Amount(5000)
	party := Contractor(0)

	// Act
	_, err = NewPremiumInfo(address, amount, party)

	// Assert
	assert.NoError(err)
}

func TestNewPremiumInfo_IncorrectAmount_Error(t *testing.T) {

	// Arrange
	assert := assert.New(t)
	address, err := btcutil.DecodeAddress(testAddress, &chaincfg.RegressionNetParams)
	assert.NoError(err)
	amount := btcutil.Amount(0)
	party := Contractor(0)

	// Act
	_, err = NewPremiumInfo(address, amount, party)

	// Assert
	assert.Error(err)
}

func TestNewPremiumInfo_IncorrectParty_Error(t *testing.T) {

	// Arrange
	assert := assert.New(t)
	address, err := btcutil.DecodeAddress(testAddress, &chaincfg.RegressionNetParams)
	assert.NoError(err)
	amount := btcutil.Amount(5000)
	party := Contractor(3)

	// Act
	_, err = NewPremiumInfo(address, amount, party)

	// Assert
	assert.Error(err)
}

func TestNewPremiumInfo_NilAddress_Error(t *testing.T) {
	// Arrange
	assert := assert.New(t)
	amount := btcutil.Amount(5000)
	party := Contractor(0)

	// Act
	_, err := NewPremiumInfo(nil, amount, party)

	// Assert
	assert.Error(err)
}

func TestNewBuilder(t *testing.T) {
	conds := newTestConditions()
	builder := NewBuilder(FirstParty, nil, NewDLC(conds))

	assert := assert.New(t)
	assert.NotNil(builder)

	assert.NotNil(builder.Contract)
}
