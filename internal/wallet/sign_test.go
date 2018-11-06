package wallet

import (
	"testing"

	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/dgarage/dlc/internal/script"
	"github.com/dgarage/dlc/internal/test"
	"github.com/stretchr/testify/assert"
)

func TestWitnessSignature(t *testing.T) {
	assert := assert.New(t)

	w, tearDownFunc := setupWallet(t)
	defer tearDownFunc()

	// pubkey and pk script
	pub, _ := w.NewPubkey()
	pkScript, _ := script.P2WPKHpkScript(pub)

	// prepare source tx
	amt := btcutil.Amount(10000)
	sourceTx := test.NewSourceTx()
	sourceTx.AddTxOut(wire.NewTxOut(int64(amt), pkScript))

	redeemTx := test.NewRedeemTx(sourceTx)

	// should fail if it's not unlocked
	_, err := w.WitnessSignature(redeemTx, 0, amt, pkScript, pub)
	assert.NotNil(err)

	// unlock for private key access
	w.Unlock(testPrivPass)

	sign, err := w.WitnessSignature(redeemTx, 0, amt, pkScript, pub)
	assert.Nil(err)

	wt := wire.TxWitness{sign, pub.SerializeCompressed()}
	redeemTx.TxIn[0].Witness = wt

	// execute script
	err = test.ExecuteScript(pkScript, redeemTx, int64(amt))
	assert.Nil(err)
}
