package dlc

import (
	"encoding/json"
	"time"

	"github.com/btcsuite/btcutil"
)

type ConditionsJSON struct {
	FixingTime     int64              `json:"fixing_time"`
	FundAmts       map[Contractor]int `json:"fund_amts"`
	FundFeerate    int                `json:"fund_feerate"`
	RedeemFeerate  int                `json:"redeem_feerate"`
	RefundLockTime uint32             `json:"refund_locktime"`
	Deals          []*DealJSON        `json:"deals"`
}

type DealJSON struct {
	Amts map[Contractor]int
	Msgs [][]byte
}

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
	amtsJson := make(map[Contractor]int)
	for c, amt := range amts {
		amtsJson[c] = int(amt)
	}

	return amtsJson
}

func dealsToJSON(deals []*Deal) []*DealJSON {
	dealsJson := []*DealJSON{}
	for _, d := range deals {
		damts := make(map[Contractor]int)
		for c, amt := range d.Amts {
			damts[c] = int(amt)
		}
		dJson := &DealJSON{Amts: damts, Msgs: d.Msgs}
		dealsJson = append(dealsJson, dJson)
	}

	return dealsJson
}

func (conds *Conditions) UnmarshalJSON(data []byte) error {
	condsJson := &ConditionsJSON{}
	err := json.Unmarshal(data, condsJson)
	if err != nil {
		return err
	}

	conds.FixingTime = time.Unix(condsJson.FixingTime, 0).UTC()

	conds.FundAmts = jsonToAmts(condsJson.FundAmts)
	conds.FundFeerate = btcutil.Amount(condsJson.FundFeerate)
	conds.RedeemFeerate = btcutil.Amount(condsJson.RedeemFeerate)
	conds.RefundLockTime = condsJson.RefundLockTime
	conds.Deals = jsonToDeals(condsJson.Deals)

	return nil
}

func jsonToAmts(amtsJson map[Contractor]int) map[Contractor]btcutil.Amount {
	amts := make(map[Contractor]btcutil.Amount)
	for c, amt := range amtsJson {
		amts[c] = btcutil.Amount(amt)
	}
	return amts
}

func jsonToDeals(dealsJson []*DealJSON) []*Deal {
	deals := []*Deal{}
	for _, dJson := range dealsJson {
		deal := &Deal{
			Amts: jsonToAmts(dJson.Amts),
			Msgs: dJson.Msgs,
		}
		deals = append(deals, deal)
	}

	return deals
}
