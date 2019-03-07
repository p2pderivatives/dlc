package wallet

import (
	"testing"

	"github.com/p2pderivatives/dlc/internal/mocks/rpcmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestNewPubkey tests generating a new public key
func TestNewPubkey(t *testing.T) {
	w, tearDownFunc := setupWallet(t)
	defer tearDownFunc()

	rpcc := &rpcmock.Client{}
	rpcc = mockImportAddress(rpcc, nil)
	w.SetRPCClient(rpcc)

	pub, err := w.NewPubkey()

	assert.Nil(t, err)
	assert.NotNil(t, pub)
}

func mockImportAddress(c *rpcmock.Client, err error) *rpcmock.Client {
	c.On("ImportAddress",
		mock.AnythingOfType("string"),
	).Return(err)

	return c
}
