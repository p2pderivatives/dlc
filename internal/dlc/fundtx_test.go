package dlc

import (
	"testing"

	"github.com/btcsuite/btcutil"
	"github.com/dgarage/dlc/internal/mocks"
	"github.com/dgarage/dlc/internal/wallet"
	"github.com/stretchr/testify/assert"
)

// Hash of block 234439
var testTxID = "14a0810ac680a3eb3f82edc878cea25ec41d6b790744e5daeef"

// setup mocke wallet
func setupTestWallet(balance btcutil.Amount) wallet.Wallet {
	w := &mocks.Wallet{}
	utxo := wallet.Utxo{
		TxID:   testTxID,
		Amount: float64(balance) / btcutil.SatoshiPerBitcoin,
	}
	w.On("ListUnspent").Return([]wallet.Utxo{utxo}, nil)
	w.On("NewWitnessPubkeyScript").Return([]byte{0x01}, nil)
	return w
}

// PrepareFundTx should fail if fund amounts aren't set
func TestPrepareFundTxNoFundAmounts(t *testing.T) {
	builder := NewBuilder(FirstParty, nil)

	err := builder.PrepareFundTxIns()
	assert.NotNil(t, err)
}

// PrepareFundTx should fail if the party doesn't have enough balance
func TestPrepareFundTxNotEnoughUtxos(t *testing.T) {
	var famt btcutil.Amount = 1 * btcutil.SatoshiPerBitcoin // 1 BTC
	var balance btcutil.Amount = famt
	testWallet := setupTestWallet(balance)
	builder := NewBuilder(FirstParty, testWallet)
	builder.SetFundAmounts(famt, famt)

	err := builder.PrepareFundTxIns()
	assert.NotNil(t, err) // not enough balance for fee
}

// PrepareFundTx should prepare the txins and txouts of fundtx
func TestPrepareFundTx(t *testing.T) {
	var balance btcutil.Amount = 1 * btcutil.SatoshiPerBitcoin // 1 BTC
	testWallet := setupTestWallet(balance)
	builder := NewBuilder(FirstParty, testWallet)

	var famt btcutil.Amount = 1 // 1 satoshi
	builder.SetFundAmounts(famt, 0)
	err := builder.PrepareFundTxIns()
	tx := builder.DLC().FundTx()

	assert := assert.New(t)
	assert.Nil(err)
	assert.NotEmpty(tx.TxIn, "tx.TxIn")
	assert.NotEmpty(tx.TxOut, "tx.TxOut")
}

// PrepareFundTx shouldn't have txouts if no changes
func TestPrepareFundTxNoChange(t *testing.T) {
	var famt btcutil.Amount = 1 * btcutil.SatoshiPerBitcoin // 1 BTC
	fee, _ := btcutil.NewAmount(float64(fundTxBaseSize + fundTxInSize))
	var balance btcutil.Amount = famt + fee
	testWallet := setupTestWallet(balance)
	builder := NewBuilder(FirstParty, testWallet)
	builder.SetFundAmounts(famt, 0)

	err := builder.PrepareFundTxIns()
	tx := builder.DLC().FundTx()

	assert := assert.New(t)
	assert.Nil(err)
	assert.NotEmpty(tx.TxIn, "tx.TxIn")
	assert.Empty(tx.TxOut, "tx.TxOut")
}
