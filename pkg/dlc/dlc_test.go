package dlc

import (
	"testing"
	"time"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/stretchr/testify/assert"
)

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
		net, ftime, famt1, famt2, frate, rrate, lc, deals)
	assert.NoError(err)

	_, err = NewConditions(
		net, time.Now(), famt1, famt2, frate, rrate, lc, deals)
	assert.Error(err)

	_, err = NewConditions(
		net, ftime, 0, famt2, frate, rrate, lc, deals)
	assert.Error(err)

	_, err = NewConditions(
		net, ftime, famt1, 0, frate, rrate, lc, deals)
	assert.Error(err)

	_, err = NewConditions(
		net, ftime, famt1, famt2, 0, rrate, lc, deals)
	assert.Error(err)

	_, err = NewConditions(
		net, ftime, famt1, famt2, frate, 0, lc, deals)
	assert.Error(err)

	_, err = NewConditions(
		net, ftime, famt1, famt2, frate, rrate, lc, []*Deal{})
	assert.Error(err)
}

func TestNewBuilder(t *testing.T) {
	builder := NewBuilder(FirstParty, nil, newTestConditions())

	assert := assert.New(t)
	assert.NotNil(builder)

	dlc := builder.DLC()
	assert.NotNil(dlc)
}
