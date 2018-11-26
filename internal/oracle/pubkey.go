package oracle

import (
	"time"

	"github.com/btcsuite/btcd/btcec"
	"github.com/dgarage/dlc/pkg/oracle"
)

// PubkeySet is an alias of oracle.PubkeySet
type PubkeySet = oracle.PubkeySet

// PubkeySet returns a key set for given fixing time
func (o *Oracle) PubkeySet(ftime time.Time) (PubkeySet, error) {
	extKey, err := o.extKeyForFixingTime(ftime)
	// derive oracle's pubkey for the given time
	if err != nil {
		return PubkeySet{}, err
	}
	pubkey, err := extKey.ECPubKey()
	if err != nil {
		return PubkeySet{}, err
	}

	// derive pubkeys for all committed R-points at the given time
	rpoints, err := committedRpoints(extKey, o.nRpoints)
	if err != nil {
		return PubkeySet{}, err
	}

	keyset := PubkeySet{
		Pubkey:           pubkey,
		CommittedRpoints: rpoints,
	}

	return keyset, nil
}

func committedRpoints(
	extKey *privExtKey, nRpoints int) ([]*btcec.PublicKey, error) {
	pubs := []*btcec.PublicKey{}
	for i := 0; i < nRpoints; i++ {
		k, err := extKey.derive(i)
		if err != nil {
			return nil, err
		}
		pub, err := k.ECPubKey()
		if err != nil {
			return nil, err
		}
		pubs = append(pubs, pub)
	}

	return pubs, nil
}
