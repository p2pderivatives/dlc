package dlc

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/btcsuite/btcutil"
)

// Deal contains information about the distributed amounts, commitment messages, and signatures of fixed messages
type Deal struct {
	Amts map[Contractor]btcutil.Amount `validate:"required,dive,gte=0"`
	Msgs [][]byte                      `validate:"required,gt=0"`
}

// NewDeal creates a new deal
func NewDeal(amt1, amt2 btcutil.Amount, msgs [][]byte) *Deal {
	amts := make(map[Contractor]btcutil.Amount)
	amts[FirstParty] = amt1
	amts[SecondParty] = amt2
	return &Deal{
		Amts: amts,
		Msgs: msgs,
	}
}

// Deal gets a deal by id
func (d *DLC) Deal(idx int) (*Deal, error) {
	if len(d.Conds.Deals) < idx+1 {
		errmsg := fmt.Sprintf("Invalid deal id. id: %d", idx)
		return nil, errors.New(errmsg)
	}

	deal := d.Conds.Deals[idx]
	return deal, nil
}

// DealByMsgs finds a deal by messages
func (d *DLC) DealByMsgs(msgs [][]byte) (idx int, deal *Deal, err error) {
	for i, deal := range d.Conds.Deals {
		if reflect.DeepEqual(deal.Msgs, msgs) {
			return i, deal, nil
		}
	}
	err = fmt.Errorf("deal not found. msgs: %v", msgs)
	return idx, deal, err
}

// FixedDealAmt returns fixed amt that the party will receive
func (b *Builder) FixedDealAmt() (btcutil.Amount, error) {
	_, deal, err := b.Contract.FixedDeal()
	if err != nil {
		return 0, err
	}

	return deal.Amts[b.party], nil
}
