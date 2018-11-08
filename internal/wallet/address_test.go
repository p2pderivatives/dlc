package wallet

import (
	"testing"

	"github.com/dgarage/dlc/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewPubkey(t *testing.T) {
	wallet, tearDownFunc := setupWallet(t)
	defer tearDownFunc()

	rpcc := &mocks.Client{}
	rpcc = mockImportAddress(rpcc, nil)
	wallet.rpc = rpcc

	pub, err := wallet.NewPubkey()

	assert.Nil(t, err)
	assert.NotNil(t, pub)
}

func mockImportAddress(c *mocks.Client, err error) *mocks.Client {
	c.On("ImportAddress",
		mock.AnythingOfType("string"),
	).Return(err)

	return c
}
