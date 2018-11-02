package dlc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBuilder(t *testing.T) {
	feeCalc := func(size int64) int64 { return size * 1 }
	builder := NewBuilder(FirstParty, nil, feeCalc)

	assert := assert.New(t)
	assert.NotNil(builder)

	dlc := builder.DLC()
	assert.NotNil(dlc)
	assert.NotNil(dlc.fundAmts, "fundAmts must exist")
	assert.NotNil(dlc.fundTxReqs, "fundTxReqs must exist")
}
