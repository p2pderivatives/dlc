package dlc

import (
	"errors"

	"github.com/btcsuite/btcd/btcec"
	"github.com/p2pderivatives/dlc/pkg/oracle"
	"github.com/p2pderivatives/dlc/pkg/schnorr"
)

// OracleRequirements contains pubkeys and commitments and signature received from oracle
type OracleRequirements struct {
	pubkeySet   *oracle.PubkeySet  // Oracle's pubkey set
	commitments []*btcec.PublicKey // Commitments for deals
	sig         []byte             // Signature for a fixed deal
	signedMsgs  [][]byte           // Messages signed by Oracle
}

func newOracleReqs(n int) *OracleRequirements {
	return &OracleRequirements{
		commitments: make([]*btcec.PublicKey, n)}
}

// PrepareOracleCommitments prepares oracle's commitments for all deals
func (d *DLC) PrepareOracleCommitments(
	V *btcec.PublicKey, Rs []*btcec.PublicKey) {
	for i, deal := range d.Conds.Deals {
		C := schnorr.CommitMulti(V, Rs, deal.Msgs)
		d.OracleReqs.commitments[i] = C
	}
}

// SetOraclePubkeySet sets oracle's pubkey set
func (b *Builder) SetOraclePubkeySet(pubset *oracle.PubkeySet) {
	b.dlc.PrepareOracleCommitments(
		pubset.Pubkey, pubset.CommittedRpoints)
	b.dlc.OracleReqs.pubkeySet = pubset
}

// FixDeal fixes a deal by setting the signature provided by oracle
func (d *DLC) FixDeal(msgs [][]byte, sigs [][]byte) error {
	dID, _, err := d.DealByMsgs(msgs)
	if err != nil {
		return err
	}

	C := d.OracleReqs.commitments[dID]
	s := schnorr.SumSigs(sigs)

	ok := schnorr.Verify(C, s)
	if !ok {
		return errors.New("invalid oracle signature")
	}

	// set fixed messages and signature for it
	d.OracleReqs.signedMsgs = msgs
	d.OracleReqs.sig = s

	return nil
}

// FixDeal fixes a deal by a oracle's signature set by picking up required messages and sigs
func (b *Builder) FixDeal(fm *oracle.SignedMsg, idxs []int) error {
	msgs := [][]byte{}
	sigs := [][]byte{}
	for _, idx := range idxs {
		msgs = append(msgs, fm.Msgs[idx])
		sigs = append(sigs, fm.Sigs[idx])
	}
	return b.dlc.FixDeal(msgs, sigs)
}

// FixedDeal returns a fixed deal
func (d *DLC) FixedDeal() (idx int, deal *Deal, err error) {
	if !d.HasDealFixed() {
		err = newNoFixedDealError()
		return
	}
	return d.DealByMsgs(d.OracleReqs.signedMsgs)
}

// HasDealFixed checks if a deal has been fixed
func (d *DLC) HasDealFixed() bool {
	return d.OracleReqs.signedMsgs != nil && d.OracleReqs.sig != nil
}
