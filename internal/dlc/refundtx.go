package dlc

import (
	"errors"

	"github.com/btcsuite/btcd/wire"
	"github.com/dgarage/dlc/internal/script"
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
	tx.LockTime = d.locktime

	// txouts
	for _, p := range []Contractor{FirstParty, SecondParty} {
		txout, err := d.refundTxOut(p)
		if err != nil {
			return nil, err
		}
		tx.AddTxOut(txout)
	}

	return tx, nil
}

func (d *DLC) refundTxOut(p Contractor) (*wire.TxOut, error) {
	famt := d.fundAmts[p]
	pub := d.fundTxReqs.pubs[p]
	pkScript, err := script.P2WPKHpkScript(pub)
	if err != nil {
		return nil, err
	}

	txout := wire.NewTxOut(int64(famt), pkScript)
	return txout, nil
}

func (b *Builder) SignRefundTx() error {
	tx, err := b.dlc.RefundTx()
	if err != nil {
		return err
	}

	sign, err := b.witsigForRedeemTx(tx)
	if err != nil {
		return err
	}

	b.dlc.refundSigns[b.party] = sign
	return nil
}

func (d *DLC) SignedRefundTx() (*wire.MsgTx, error) {
	tx, err := d.RefundTx()
	if err != nil {
		return nil, err
	}

	wt, err := d.witessForRefundTx()
	if err != nil {
		return nil, err
	}

	tx.TxIn[0].Witness = wt
	return tx, nil
}

func (d *DLC) witessForRefundTx() (wire.TxWitness, error) {
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
