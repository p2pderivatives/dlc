package oracle

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPubkeySet(t *testing.T) {
	assert := assert.New(t)

	o, _ := newTestOracle(t)

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
