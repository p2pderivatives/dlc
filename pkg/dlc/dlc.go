package dlc

import (
	"bytes"
	"encoding/hex"
	"errors"
	"time"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/dgarage/dlc/pkg/script"
	"github.com/dgarage/dlc/pkg/wallet"
	validator "gopkg.in/go-playground/validator.v9"
)

// DLC contains all information required for DLC contract
// including FundTx, SettlementTx, RefundTx
type DLC struct {
	Conds *Conditions

	// requirements
	Pubs       map[Contractor]*btcec.PublicKey // pubkeys used for script and txout
	FundTxReqs *FundTxRequirements             // fund txins/outs
	OracleReqs *OracleRequirements
	RefundSigs map[Contractor][]byte // signatures for refund tx
	ExecSigs   [][]byte              // counterparty's signatures for CETxs
}

// NewDLC initializes DLC
func NewDLC(conds *Conditions) *DLC {
	nDeal := len(conds.Deals)
	return &DLC{
		Conds:      conds,
		Pubs:       make(map[Contractor]*btcec.PublicKey),
		FundTxReqs: NewFundTxReqs(),
		OracleReqs: newOracleReqs(nDeal),
		RefundSigs: make(map[Contractor][]byte),
		ExecSigs:   make([][]byte, nDeal),
	}
}

// Conditions contains conditions of a contract
type Conditions struct {
	FixingTime     time.Time                     `validate:"required,gt=time.Now()"`
	FundAmts       map[Contractor]btcutil.Amount `validate:"required,dive,gt=0"`
	FundFeerate    btcutil.Amount                `validate:"required,gt=0"` // fund fee rate (satoshi per byte)
	RedeemFeerate  btcutil.Amount                `validate:"required,gt=0"` // redeem fee rate (satoshi per byte)
	RefundLockTime uint32                        `validate:"required,gt=0"` // refund locktime (block height)
	Deals          []*Deal                       `validate:"required,gt=0,dive,required"`
}

// NewConditions creates a new DLC conditions
func NewConditions(
	ftime time.Time,
	famt1, famt2 btcutil.Amount,
	ffeerate, rfeerate btcutil.Amount, // fund feerate and redeem feerate
	refundLockTime uint32, // refund locktime
	deals []*Deal,
) (*Conditions, error) {
	famts := make(map[Contractor]btcutil.Amount)
	famts[FirstParty] = famt1
	famts[SecondParty] = famt2

	conds := &Conditions{
		FixingTime:     ftime,
		FundAmts:       famts,
		FundFeerate:    ffeerate,
		RedeemFeerate:  rfeerate,
		RefundLockTime: refundLockTime,
		Deals:          deals,
	}

	// validate structure
	err := validator.New().Struct(conds)

	return conds, err
}

// ClosingTxOut returns a final txout owned only by a given party
func (d *DLC) ClosingTxOut(
	p Contractor, amt btcutil.Amount) (*wire.TxOut, error) {
	pub := d.Pubs[p]
	if pub == nil {
		return nil, errors.New("missing pubkey")
	}

	pkScript, err := script.P2WPKHpkScript(pub)
	if err != nil {
		return nil, err
	}

	txout := wire.NewTxOut(int64(amt), pkScript)
	return txout, nil
}

const txVersion = 2

// Contractor represents a contractor type
type Contractor int

const (
	// FirstParty is a contractor who creates offer
	FirstParty Contractor = 0
	// SecondParty is a contractor who accepts offer
	SecondParty Contractor = 1
)

// counterparty returns the counterparty
func counterparty(p Contractor) (cp Contractor) {
	switch p {
	case FirstParty:
		cp = SecondParty
	case SecondParty:
		cp = FirstParty
	}
	return cp
}

// Builder builds DLC by interacting with wallet
type Builder struct {
	party  Contractor
	wallet wallet.Wallet
	dlc    *DLC
}

// NewBuilder creates a new Builder for a contractor
func NewBuilder(
	p Contractor, w wallet.Wallet, conds *Conditions) *Builder {
	return &Builder{
		dlc:    NewDLC(conds),
		party:  p,
		wallet: w,
	}
}

// DLC returns the DLC constructed by builder
func (b *Builder) DLC() *DLC {
	return b.dlc
}

// PreparePubkey sets fund pubkey
func (b *Builder) PreparePubkey() error {
	pub, err := b.wallet.NewPubkey()
	if err != nil {
		return err
	}
	b.dlc.Pubs[b.party] = pub
	return nil
}

// CopyReqsFromCounterparty copies requirements from counterparty
func (b *Builder) CopyReqsFromCounterparty(d *DLC) {
	p := counterparty(b.party)

	// pubkey
	b.dlc.Pubs[p] = d.Pubs[p]

	// fund requirements
	fundReqs := d.FundTxReqs
	b.dlc.FundTxReqs.TxIns[p] = fundReqs.TxIns[p]
	b.dlc.FundTxReqs.TxOut[p] = fundReqs.TxOut[p]
}

func txToHex(tx *wire.MsgTx) (string, error) {
	// Serialize the transaction and convert to hex string.
	buf := bytes.NewBuffer(make([]byte, 0, tx.SerializeSize()))
	if err := tx.Serialize(buf); err != nil {
		return "", err
	}
	h := hex.EncodeToString(buf.Bytes())
	return h, nil
}

func hexToTx(txHex string) (tx *wire.MsgTx, err error) {
	txbin, err := hex.DecodeString(txHex)
	if err != nil {
		return nil, err
	}
	bufr := bytes.NewReader(txbin)
	err = tx.Deserialize(bufr)
	return
}
