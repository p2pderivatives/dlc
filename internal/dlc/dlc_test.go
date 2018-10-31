package dlc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateDraft(t *testing.T) {
	feeCalc := func(size int64) int64 { return size * 1 }
	builder := NewBuilder(FirstParty, nil, feeCalc)
	assert.NotNil(t, builder.DLC())
}
