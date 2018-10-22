package oracle

import (
	"math/big"
	"testing"
	"time"

	"github.com/btcsuite/btcd/btcec"
	"github.com/stretchr/testify/assert"
)

func TestPubkeySet(t *testing.T) {
	assert := assert.New(t)

	o, _ := newTestOracle(t, 1)

	// Get KeySet
	ftime := time.Now()
	keyset, err := o.PubkeySet(ftime)
	assert.Nil(err)
	assert.IsType("", keyset.Pubkey)
	assert.IsType([]string{}, keyset.CommittedRpoints)

	// Compare with other keysets
	keysetSame, _ := o.PubkeySet(ftime) // same time
	assert.Equal(keyset, keysetSame)

	keysetNextYear, _ := o.PubkeySet(ftime.AddDate(1, 0, 0)) // next year
	assert.NotEqual(keyset, keysetNextYear)

	keysetNextMonth, _ := o.PubkeySet(ftime.AddDate(0, 1, 0)) // next month
	assert.NotEqual(keyset, keysetNextMonth)

	keysetTomorrow, _ := o.PubkeySet(ftime.AddDate(0, 0, 1)) // tomorrow
	assert.NotEqual(keyset, keysetTomorrow)

	keysetHourLater, _ := o.PubkeySet(ftime.Add(1 * time.Hour)) // an hour later
	assert.NotEqual(keyset, keysetHourLater)

	keysetMiniteLater, _ := o.PubkeySet(ftime.Add(1 * time.Minute)) // a minute later
	assert.NotEqual(keyset, keysetMiniteLater)

	keysetSecondLater, _ := o.PubkeySet(ftime.Add(1 * time.Second)) // a second later
	assert.NotEqual(keyset, keysetSecondLater)
}

func TestCommit(t *testing.T) {
	assert := assert.New(t)

	// Get PubkeySet
	o, _ := newTestOracle(t, 6)
	ftime := time.Now().AddDate(0, 0, 1)
	keyset, _ := o.PubkeySet(ftime)

	O, _ := StrToPubkey(keyset.Pubkey)
	Rs := keyset.CommittedRpoints

	vals := []int{1, 2, 3, 4, 5, 6} // commiting value 123456
	mkey := new(btcec.PublicKey)    // concatenated pub key

	for i, val := range vals {
		R, _ := StrToPubkey(Rs[i])
		m := big.NewInt(int64(val)).Bytes()
		// R is contract key,O is oracle public key.
		// R - H(R,m)O
		p := o.Commit(R, O, m)

		// If there are multiple messages, concatenate public keys.
		if mkey.X == nil {
			mkey.X, mkey.Y = p.X, p.Y
		} else {
			mkey.X, mkey.Y = btcec.S256().Add(mkey.X, mkey.Y, mkey.X, mkey.Y)
		}
	}

	// TODO: add meanigful assertions
	assert.NotNil(mkey)
}
