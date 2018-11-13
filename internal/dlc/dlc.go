package dlc

import (
	"errors"
	"fmt"

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

	// requirements to execute DLC
	pubs        map[Contractor]*btcec.PublicKey
	fundTxReqs  *FundTxRequirements
	refundSigns map[Contractor][]byte
}

func newDLC(conds Conditions) *DLC {
	return &DLC{
		conds:       conds,
		pubs:        make(map[Contractor]*btcec.PublicKey),
		fundTxReqs:  newFundTxReqs(),
		refundSigns: make(map[Contractor][]byte),
	}
}

// Conditions contains conditions of a contract
type Conditions struct {
	FundAmts      map[Contractor]btcutil.Amount `validate:"required,dive,gt=0"`
	FundFeerate   btcutil.Amount                `validate:"required,gt=0"` // fund fee rate (satoshi per byte)
	RedeemFeerate btcutil.Amount                `validate:"required,gt=0"` // redeem fee rate (satoshi per byte)
	LockTime      uint32                        `validate:"required,gt=0"` // refund locktime (block height)
}

// NewConditions creates a new DLC conditions
func NewConditions(
	famt1, famt2 btcutil.Amount,
	ffeerate, rfeerate btcutil.Amount, // fund feerate and redeem feerate
	lc uint32, // locktime
) (Conditions, error) {
	famts := make(map[Contractor]btcutil.Amount)
	famts[FirstParty] = famt1
	famts[SecondParty] = famt2

	conds := Conditions{
		FundAmts:      famts,
		FundFeerate:   ffeerate,
		RedeemFeerate: rfeerate,
		LockTime:      lc,
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
func (d *DLC) counterparty(p Contractor) (cp Contractor) {
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
	p := b.dlc.counterparty(b.party)

	// pubkey
	b.dlc.pubs[p] = d.pubs[p]

	// fund requirements
	fundReqs := d.fundTxReqs
	b.dlc.fundTxReqs.txIns[p] = fundReqs.txIns[p]
	b.dlc.fundTxReqs.txOut[p] = fundReqs.txOut[p]
}

// AcceptCounterpartySign verifies couterparty's given sign is valid and then
func (b *Builder) AcceptCounterpartySign(sign []byte) error {
	p := b.dlc.counterparty(b.party)

	err := b.dlc.VerifyRefundTx(sign, b.dlc.pubs[p])
	if err != nil {
		return fmt.Errorf("counterparty's signature didn't pass verification, had error: %v", err)
	}

	// sign passed verification, accept it
	b.dlc.refundSigns[p] = sign
	return nil
}
