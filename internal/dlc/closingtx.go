package dlc

import (
	"math/big"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/dgarage/dlc/internal/script"
	"github.com/dgarage/dlc/internal/wallet"
)

// closingTxSize is size of closing tx
const closingTxSize = 216

// closingTxOutAt is a txout index of contract execution tx
const closingTxOutAt = 0

// ClosingTx constructs a tx that redeems a given CET
func (d *DLC) ClosingTx(
	p Contractor, cetx *wire.MsgTx) (*wire.MsgTx, error) {

	tx := wire.NewMsgTx(txVersion)

	// txin
	txid := cetx.TxHash()
	txin := wire.NewTxIn(
		wire.NewOutPoint(&txid, closingTxOutAt), nil, nil)

	tx.AddTxIn(txin)

	// TODO: the party who sends closing tx have to pay for this fee.
	// Consider spliting it with the counterparty
	in := btcutil.Amount(cetx.TxOut[closingTxOutAt].Value)
	fee := d.redeemTxFee(closingTxSize)
	out := in - fee

	if out <= 0 {
		return nil, newNotEnoughFeesError(in, fee)
	}

	txout, err := d.ClosingTxOut(p, out)
	if err != nil {
		return nil, err
	}
	tx.AddTxOut(txout)

	return tx, nil
}

// SignedClosingTx constructs a closing tx with witness
func (b *Builder) SignedClosingTx(cetx *wire.MsgTx) (*wire.MsgTx, error) {
	dID, _, err := b.dlc.FixedDeal()
	if err != nil {
		return nil, err
	}
	C := b.dlc.oracleReqs.commitments[dID]

	tx, err := b.dlc.ClosingTx(b.party, cetx)
	if err != nil {
		return nil, err
	}

	wit, err := b.witnessForCEScript(tx, cetx, C)
	if err != nil {
		return nil, err
	}
	tx.TxIn[0].Witness = wit

	return tx, nil
}

func (b *Builder) witnessForCEScript(
	tx, cetx *wire.MsgTx, C *btcec.PublicKey) (wire.TxWitness, error) {
	cetxout := cetx.TxOut[closingTxOutAt]
	amt := btcutil.Amount(cetxout.Value)

	cparty := counterparty(b.party)
	pub1, pub2 := b.dlc.pubs[b.party], b.dlc.pubs[cparty]

	sc, err := script.ContractExecutionScript(
		pub1, pub2, C)
	if err != nil {
		return nil, err
	}

	// callback function that adds message sign to private key
	osign := b.dlc.oracleReqs.sign
	privkeyConverter := genAddSignToPrivkeyFunc(osign)

	sign, err := b.wallet.WitnessSignatureWithCallback(
		tx, closingTxOutAt, amt, sc, pub1, privkeyConverter)
	if err != nil {
		return nil, err
	}

	wit := script.WitnessForCEScript(sign, sc)
	return wit, nil
}

func genAddSignToPrivkeyFunc(
	sign []byte) wallet.PrivateKeyConverter {
	return func(priv *btcec.PrivateKey) (*btcec.PrivateKey, error) {
		n := new(big.Int).Add(priv.D, new(big.Int).SetBytes(sign))
		n = new(big.Int).Mod(n, btcec.S256().N)
		p, _ := btcec.PrivKeyFromBytes(btcec.S256(), n.Bytes())
		return p, nil
	}
}
