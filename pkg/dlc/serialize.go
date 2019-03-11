package dlc

import (
	"encoding/json"
	"reflect"
	"time"

	"github.com/btcsuite/btcutil"
	"github.com/p2pderivatives/dlc/pkg/utils"
)

// ConditionsJSON is contract conditions in JSON format
type ConditionsJSON struct {
	FixingTime     int64              `json:"fixing_time"`
	FundAmts       map[Contractor]int `json:"fund_amts"`
	FundFeerate    int                `json:"fund_feerate"`
	RedeemFeerate  int                `json:"redeem_feerate"`
	RefundLockTime uint32             `json:"refund_locktime"`
	Deals          []*DealJSON        `json:"deals"`
}

// DealJSON is DLC deals in JSON format
type DealJSON struct {
	Amts map[Contractor]int `json:"amts"`
	Msgs [][]byte           `json:"msgs"`
}

// PublicKeys is public keys in hex string format
type PublicKeys map[Contractor]string

// Addresses is addresses in string
type Addresses map[Contractor]string

// MarshalJSON implements json.Marshaler
func (conds *Conditions) MarshalJSON() ([]byte, error) {
	return json.Marshal(&ConditionsJSON{
		FixingTime:     conds.FixingTime.Unix(),
		FundAmts:       amtsToJSON(conds.FundAmts),
		FundFeerate:    int(conds.FundFeerate),
		RedeemFeerate:  int(conds.RedeemFeerate),
		RefundLockTime: conds.RefundLockTime,
		Deals:          dealsToJSON(conds.Deals),
	})
}

func amtsToJSON(amts map[Contractor]btcutil.Amount) map[Contractor]int {
	amtsJSON := make(map[Contractor]int)
	for c, amt := range amts {
		amtsJSON[c] = int(amt)
	}

	return amtsJSON
}

func dealsToJSON(deals []*Deal) []*DealJSON {
	dealsJSON := []*DealJSON{}
	for _, d := range deals {
		damts := make(map[Contractor]int)
		for c, amt := range d.Amts {
			damts[c] = int(amt)
		}
		dJSON := &DealJSON{Amts: damts, Msgs: d.Msgs}
		dealsJSON = append(dealsJSON, dJSON)
	}

	return dealsJSON
}

// UnmarshalJSON implements json.Unmarshaler
func (conds *Conditions) UnmarshalJSON(data []byte) error {
	condsJSON := &ConditionsJSON{}
	err := json.Unmarshal(data, condsJSON)
	if err != nil {
		return err
	}

	conds.FixingTime = time.Unix(condsJSON.FixingTime, 0).UTC()

	conds.FundAmts = jsonToAmts(condsJSON.FundAmts)
	conds.FundFeerate = btcutil.Amount(condsJSON.FundFeerate)
	conds.RedeemFeerate = btcutil.Amount(condsJSON.RedeemFeerate)
	conds.RefundLockTime = condsJSON.RefundLockTime
	conds.Deals = jsonToDeals(condsJSON.Deals)

	return nil
}

func jsonToAmts(amtsJSON map[Contractor]int) map[Contractor]btcutil.Amount {
	amts := make(map[Contractor]btcutil.Amount)
	for c, amt := range amtsJSON {
		amts[c] = btcutil.Amount(amt)
	}
	return amts
}

func jsonToDeals(dealsJSON []*DealJSON) []*Deal {
	deals := []*Deal{}
	for _, dJSON := range dealsJSON {
		deal := &Deal{
			Amts: jsonToAmts(dJSON.Amts),
			Msgs: dJSON.Msgs,
		}
		deals = append(deals, deal)
	}

	return deals
}

// PublicKeys converts btcec.PublicKey to hex string
func (d *DLC) PublicKeys() PublicKeys {
	pubs := make(PublicKeys)
	for c, p := range d.Pubs {
		pubs[c] = utils.PubkeyToStr(p)
	}
	return pubs
}

// ParsePublicKeys parses public key hex strings
func (d *DLC) ParsePublicKeys(pubs PublicKeys) error {
	for c, p := range pubs {
		pub, err := utils.ParsePublicKey(p)
		if err != nil {
			return err
		}
		d.Pubs[c] = pub
	}
	return nil
}

// Addresses converts btcutil.Address to string
func (d *DLC) Addresses() Addresses {
	addrs := make(Addresses)
	for c, addr := range d.Addrs {
		if reflect.ValueOf(addr).IsNil() != true {
			addrs[c] = addr.EncodeAddress()
		}
	}
	return addrs
}

// ParseAddresses parses address string
func (d *DLC) ParseAddresses(addrs Addresses) error {
	for c, addrStr := range addrs {
		addr, err := btcutil.DecodeAddress(addrStr, d.NetParams)
		if err != nil {
			return err
		}
		d.Addrs[c] = addr
	}
	return nil
}

// ChangeAddresses converts btcutil.Address to string
func (d *DLC) ChangeAddresses() Addresses {
	addrs := make(Addresses)
	for c, addr := range d.ChangeAddrs {
		if reflect.ValueOf(addr).IsNil() != true {
			addrs[c] = addr.EncodeAddress()
		}
	}
	return addrs
}

// ParseChangeAddresses parses address string
func (d *DLC) ParseChangeAddresses(addrs Addresses) error {
	for c, addrStr := range addrs {
		addr, err := btcutil.DecodeAddress(addrStr, d.NetParams)
		if err != nil {
			return err
		}
		d.ChangeAddrs[c] = addr
	}
	return nil
}
