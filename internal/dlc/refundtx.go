package dlc

import (
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

// RefundTx creates refund tx
// refund transaction
// input:
//   [0]:fund transaction output[0]
//       Sequence (0xfeffffff LE)
// output:
//   [0]:p2wpkh a
//   [1]:p2wpkh b
// locktime:
//    Value decided by contract.
func (d *DLC) RefundTx() (*wire.MsgTx, error) {
	tx, err := d.newRedeemTx()
	if err != nil {
		return nil, err
	}

	// use locktime
	tx.TxIn[fundTxInAt].Sequence-- // max(0xffffffff-0x01)
	tx.LockTime = d.lockTime

	// txouts
	for _, p := range []Contractor{FirstParty, SecondParty} {
		txout, err := d.ClosingTxOut(p, d.fundAmts[p])

		if err != nil {
			fmt.Printf("err in closing tx out:   %+v\n", err)
			return nil, err
		}
		tx.AddTxOut(txout)
	}

	return tx, nil
}

// SignRefundTx creates signature for a refund tx, sets it, and returns it
func (b *Builder) SignRefundTx() ([]byte, error) {
	tx, err := b.dlc.RefundTx()
	if err != nil {
		return nil, err
	}

	sign, err := b.witsigForFundTxIn(tx)
	if err != nil {
		return nil, err
	}

	b.dlc.refundSigns[b.party] = sign

	return sign, nil
}

// SignedRefundTx returns a refund tx with its witness signature
func (d *DLC) SignedRefundTx() (*wire.MsgTx, error) {
	tx, err := d.RefundTx()
	if err != nil {
		return nil, err
	}

	wt, err := d.witnessForRefundTx()
	if err != nil {
		return nil, err
	}

	tx.TxIn[0].Witness = wt
	return tx, nil
}

func (d *DLC) witnessForRefundTx() (wire.TxWitness, error) {
	sc, err := d.fundScript()
	if err != nil {
		return nil, err
	}

	sign1 := d.refundSigns[FirstParty]
	if sign1 == nil {
		return nil, errors.New("First party must sign refund tx")
	}

	sign2 := d.refundSigns[SecondParty]
	if sign2 == nil {
		return nil, errors.New("Second party must sign refund tx")
	}

	wt := wire.TxWitness{[]byte{}, sign1, sign2, sc}
	return wt, nil
}

// VerifyRefundTx verifies the refund transaction signature.
// Returns nil if the passed in sign is valid and corresponds to the passed
// in public key, and an error if it isnt.
func (d *DLC) VerifyRefundTx(sign []byte, pub *btcec.PublicKey) error {
	// parse signature
	s, err := btcec.ParseDERSignature(sign, btcec.S256())
	if err != nil {
		return err
	}

	// verify
	script, err := d.fundScript()
	if err != nil {
		return err
	}
	if script == nil {
		return fmt.Errorf("fund script not found ")
	}
	tx, err := d.RefundTx()
	if err != nil {
		return err
	}
	sighashes := txscript.NewTxSigHashes(tx)

	fundtx, err := d.FundTx()
	if err != nil {
		return err
	}
	amt := fundtx.TxOut[0].Value

	hash, err := txscript.CalcWitnessSigHash(script, sighashes, txscript.SigHashAll,
		tx, 0, amt)
	if err != nil {
		return err
	}

	verify := s.Verify(hash, pub)
	if !verify {
		return fmt.Errorf("verify fail : %v", verify)
	}

	return nil
}
