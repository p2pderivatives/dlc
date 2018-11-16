package dlc

import (
	"testing"

	"github.com/btcsuite/btcd/btcec"
	"github.com/dgarage/dlc/internal/oracle"
	"github.com/dgarage/dlc/internal/test"
	"github.com/stretchr/testify/assert"
)

func TestSetOraclePubkeySet(t *testing.T) {
	b, _, dID := setupContractorForOracleTest()

	_, pub := test.RandKeys()
	_, R := test.RandKeys()
	pubset := &oracle.PubkeySet{
		Pubkey: pub, CommittedRpoints: []*btcec.PublicKey{R}}

	b.SetOraclePubkeySet(pubset)
	d := b.DLC()

	assert.NotNil(t, d.oracleReqs.commitments[dID])
}

func TestFixDeal(t *testing.T) {
	assert := assert.New(t)
	var err error

	// setup
	b, deal, dID := setupContractorForOracleTest()
	privkey, C := test.RandKeys()
	b.dlc.oracleReqs.commitments[dID] = C

	// fail with invalid sign
	privInvalid, _ := test.RandKeys()
	osignsInvalid := [][]byte{privInvalid.D.Bytes()}
	osignsetInvalid := &oracle.SignSet{Msgs: deal.Msgs, Signs: osignsInvalid}

	err = b.FixDeal(osignsetInvalid, []int{0})
	assert.Error(err)

	// success with valid sign and message set
	osigns := [][]byte{privkey.D.Bytes()}
	osignset := &oracle.SignSet{Msgs: deal.Msgs, Signs: osigns}
	err = b.FixDeal(osignset, []int{0})
	assert.NoError(err)

	// retrieve fixed deal
	fixedID, fixedDeal, err := b.dlc.FixedDeal()
	assert.NoError(err)
	assert.Equal(dID, fixedID)
	assert.Equal(deal, fixedDeal)
}

func setupContractorForOracleTest() (*Builder, *Deal, int) {
	conds := newTestConditions()

	// set deals
	msgs := [][]byte{{1}}
	deal := NewDeal(1, 1, msgs)
	conds.Deals = []*Deal{deal}

	// init first party
	w := setupTestWallet()
	w = mockSelectUnspent(w, 1, 1, nil)
	b := NewBuilder(FirstParty, w, conds)

	dID, _, _ := b.dlc.DealByMsgs(msgs)

	return b, deal, dID
}
