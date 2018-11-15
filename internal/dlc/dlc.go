package dlc

import (
	"errors"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/dgarage/dlc/internal/script"
	"github.com/dgarage/dlc/internal/wallet"
	validator "gopkg.in/go-playground/validator.v9"
)

// DLC contains all information required for DLC contract
// including FundTx, SettlementTx, RefundTx
type DLC struct {
	conds Conditions

	// requirements
	pubs        map[Contractor]*btcec.PublicKey // pubkeys used for script and txout
	fundTxReqs  *FundTxRequirements             // fund txins/outs
	oracleReqs  *OracleRequirements
	refundSigns map[Contractor][]byte // counterparty's sign for refund tx
	cetxSigns   [][]byte              // counterparty's signs for CETs
}

func newDLC(conds Conditions) *DLC {
	nDeal := len(conds.Deals)
	return &DLC{
		conds:       conds,
		pubs:        make(map[Contractor]*btcec.PublicKey),
		fundTxReqs:  newFundTxReqs(),
		oracleReqs:  newOracleReqs(nDeal),
		refundSigns: make(map[Contractor][]byte),
		cetxSigns:   make([][]byte, nDeal),
	}
}

// Conditions contains conditions of a contract
type Conditions struct {
	FundAmts      map[Contractor]btcutil.Amount `validate:"required,dive,gt=0"`
	FundFeerate   btcutil.Amount                `validate:"required,gt=0"` // fund fee rate (satoshi per byte)
	RedeemFeerate btcutil.Amount                `validate:"required,gt=0"` // redeem fee rate (satoshi per byte)
	// TODO: add SettlementAt
	// SettlementAt  time.Time                     `validate:"required"`
	LockTime uint32  `validate:"required,gt=0"` // refund locktime (block height)
	Deals    []*Deal `validate:"required,gt=0,dive,required"`
}

// NewConditions creates a new DLC conditions
func NewConditions(
	famt1, famt2 btcutil.Amount,
	ffeerate, rfeerate btcutil.Amount, // fund feerate and redeem feerate
	lc uint32, // locktime
	deals []*Deal,
) (Conditions, error) {
	famts := make(map[Contractor]btcutil.Amount)
	famts[FirstParty] = famt1
	famts[SecondParty] = famt2

	conds := Conditions{
		FundAmts:      famts,
		FundFeerate:   ffeerate,
		RedeemFeerate: rfeerate,
		LockTime:      lc,
		Deals:         deals,
	}

	// validate structure
	err := validator.New().Struct(conds)

	return conds, err
}

// ClosingTxOut returns a final txout owned only by a given party
func (d *DLC) ClosingTxOut(
	p Contractor, amt btcutil.Amount) (*wire.TxOut, error) {
	pub := d.pubs[p]
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
	p Contractor, w wallet.Wallet, conds Conditions) *Builder {
	return &Builder{
		dlc:    newDLC(conds),
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
	b.dlc.pubs[b.party] = pub
	return nil
}

// CopyReqsFromCounterparty copies requirements from counterparty
func (b *Builder) CopyReqsFromCounterparty(d *DLC) {
	p := counterparty(b.party)

	// pubkey
	b.dlc.pubs[p] = d.pubs[p]

	// fund requirements
	fundReqs := d.fundTxReqs
	b.dlc.fundTxReqs.txIns[p] = fundReqs.txIns[p]
	b.dlc.fundTxReqs.txOut[p] = fundReqs.txOut[p]
}
