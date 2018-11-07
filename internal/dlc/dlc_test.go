package dlc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBuilder(t *testing.T) {
	builder := NewBuilder(FirstParty, nil)

	assert := assert.New(t)
	assert.NotNil(builder)

	dlc := builder.DLC()
	assert.NotNil(dlc)
	assert.NotNil(dlc.fundAmts, "fundAmts must exist")
	assert.NotNil(dlc.fundTxReqs, "fundTxReqs must exist")
}
