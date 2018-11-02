package dlc

import (
	"testing"

	"github.com/btcsuite/btcutil"
	"github.com/stretchr/testify/assert"
)

var testFeePerByte btcutil.Amount = 1
var testFeeCalc = func(size int64) btcutil.Amount {
	return testFeePerByte.MulF64(float64(size))
}

func TestNewBuilder(t *testing.T) {
	builder := NewBuilder(FirstParty, nil, testFeeCalc)

	assert := assert.New(t)
	assert.NotNil(builder)

	dlc := builder.DLC()
	assert.NotNil(dlc)
	assert.NotNil(dlc.fundAmts, "fundAmts must exist")
	assert.NotNil(dlc.fundTxReqs, "fundTxReqs must exist")
}
