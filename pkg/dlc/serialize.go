package dlc

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/p2pderivatives/dlc/pkg/oracle"
	"github.com/p2pderivatives/dlc/pkg/utils"
)

// OracleJSON is oracle information in JSON format
type OracleJSON struct {
	PubkeySet   *oracle.PubkeySetJSON `json:"pubkey"`
	Commitments []string              `json:"commitments"`
	Sig         []byte                `json:"sig"`
	SignedMsgs  [][]byte              `json:"signed_msgs"`
}

// ConditionsJSON is contract conditions in JSON format
type ConditionsJSON struct {
	Net            string             `json:"network"`
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
func (o *Oracle) MarshalJSON() ([]byte, error) {
	var pubkeyJSON *oracle.PubkeySetJSON
	if o.PubkeySet != nil {
		pubkeyJSON = o.PubkeySet.JSON()
	}

	Cs := []string{}
	for _, c := range o.Commitments {
		Cs = append(Cs, utils.PubkeyToStr(c))
	}

	return json.Marshal(&OracleJSON{
		PubkeySet:   pubkeyJSON,
		Commitments: Cs,
		Sig:         o.Sig,
		SignedMsgs:  o.SignedMsgs,
	})
}

// MarshalJSON implements json.Marshaler
func (conds *Conditions) MarshalJSON() ([]byte, error) {
	return json.Marshal(&ConditionsJSON{
		Net:            conds.NetParams.Name,
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
func (o *Oracle) UnmarshalJSON(data []byte) error {
	oJSON := &OracleJSON{}
	err := json.Unmarshal(data, oJSON)
	if err != nil {
		return err
	}

	if oJSON.PubkeySet != nil {
		pubset := &oracle.PubkeySet{}
		err = pubset.ParseJSON(oJSON.PubkeySet)
		if err != nil {
			return err
		}
		o.PubkeySet = pubset
	}

	for k, cstr := range oJSON.Commitments {
		c, err := utils.ParsePublicKey(cstr)
		if err != nil {
			return err
		}
		o.Commitments[k] = c
	}

	o.Sig = oJSON.Sig
	o.SignedMsgs = oJSON.SignedMsgs

	return nil
}

// UnmarshalJSON implements json.Unmarshaler
func (conds *Conditions) UnmarshalJSON(data []byte) error {
	condsJSON := &ConditionsJSON{}
	err := json.Unmarshal(data, condsJSON)
	if err != nil {
		return err
	}

	net, err := strToNetParams(condsJSON.Net)
	if err != nil {
		return err
	}
	conds.NetParams = net

	conds.FixingTime = time.Unix(condsJSON.FixingTime, 0).UTC()

	conds.FundAmts = jsonToAmts(condsJSON.FundAmts)
	conds.FundFeerate = btcutil.Amount(condsJSON.FundFeerate)
	conds.RedeemFeerate = btcutil.Amount(condsJSON.RedeemFeerate)
	conds.RefundLockTime = condsJSON.RefundLockTime
	conds.Deals = jsonToDeals(condsJSON.Deals)

	return nil
}

// InvalidNetworkNameError is used when invalid network name is given
type InvalidNetworkNameError struct{ error }

func strToNetParams(str string) (*chaincfg.Params, error) {
	var net *chaincfg.Params
	var err error
	switch str {
	case chaincfg.MainNetParams.Name:
		net = &chaincfg.MainNetParams
	case chaincfg.TestNet3Params.Name:
		net = &chaincfg.TestNet3Params
	case chaincfg.RegressionNetParams.Name:
		net = &chaincfg.RegressionNetParams
	case chaincfg.SimNetParams.Name:
		net = &chaincfg.SimNetParams
	default:
		msg := fmt.Errorf("invalid network name. %s", str)
		err = InvalidNetworkNameError{error: msg}
	}

	return net, err
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
		if !reflect.ValueOf(addr).IsNil() {
			addrs[c] = addr.EncodeAddress()
		}
	}
	return addrs
}

// ParseAddresses parses address string
func (d *DLC) ParseAddresses(addrs Addresses) error {
	for c, addrStr := range addrs {
		addr, err := btcutil.DecodeAddress(addrStr, d.Conds.NetParams)
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
		if !reflect.ValueOf(addr).IsNil() {
			addrs[c] = addr.EncodeAddress()
		}
	}
	return addrs
}

// ParseChangeAddresses parses address string
func (d *DLC) ParseChangeAddresses(addrs Addresses) error {
	for c, addrStr := range addrs {
		addr, err := btcutil.DecodeAddress(addrStr, d.Conds.NetParams)
		if err != nil {
			return err
		}
		d.ChangeAddrs[c] = addr
	}
	return nil
}
