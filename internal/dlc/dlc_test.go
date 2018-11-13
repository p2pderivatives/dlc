package dlc

import (
	"testing"

	"github.com/btcsuite/btcutil"
	"github.com/stretchr/testify/assert"
)

func TestCondions(t *testing.T) {
	assert := assert.New(t)

	var famt1, famt2,
		frate, rrate btcutil.Amount = 1, 1, 1, 1
	var lc uint32 = 1

	var err error
	_, err = NewConditions(famt1, famt2, frate, rrate, lc)
	assert.Nil(err)

	_, err = NewConditions(0, famt2, frate, rrate, lc)
	assert.NotNil(err)

	_, err = NewConditions(famt1, 0, frate, rrate, lc)
	assert.NotNil(err)

	_, err = NewConditions(famt1, famt2, 0, rrate, lc)
	assert.NotNil(err)

	_, err = NewConditions(famt1, famt2, frate, 0, lc)
	assert.NotNil(err)

	_, err = NewConditions(famt1, famt2, frate, rrate, 0)
	assert.NotNil(err)
}

func TestNewBuilder(t *testing.T) {
	conds, _ := NewConditions(1, 1, 1, 1, 1)
	builder := NewBuilder(FirstParty, nil, conds)

	assert := assert.New(t)
	assert.NotNil(builder)

	dlc := builder.DLC()
	assert.NotNil(dlc)
	assert.NotNil(dlc.fundTxReqs, "fundTxReqs must exist")
}
