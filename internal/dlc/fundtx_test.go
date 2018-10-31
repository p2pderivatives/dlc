package dlc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrepareFundTx(t *testing.T) {
	assert := assert.New(t)

	feeCalc := func(size int64) int64 { return size * 1 }
	builder := NewBuilder(FirstParty, nil, feeCalc)

	// fail if fund amounts aren't set
	err := builder.PrepareFundTx()
	assert.NotNil(err)

	builder.SetFundAmounts(1, 1)
	err = builder.PrepareFundTx()
	assert.Nil(err)

	dlc := builder.DLC()
	tx := dlc.FundTx()
	assert.NotEmpty(tx.TxIn)
	assert.NotEmpty(tx.TxOut)
}
