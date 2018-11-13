package dlc

import (
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/dgarage/dlc/internal/mocks"
	"github.com/dgarage/dlc/internal/script"
	"github.com/dgarage/dlc/internal/test"
	"github.com/dgarage/dlc/internal/wallet"
	"github.com/stretchr/testify/mock"
)

// setup mocke wallet
func setupTestWallet() *mocks.Wallet {
	w := &mocks.Wallet{}
	priv, pub := test.RandKeys()
	w.On("NewPubkey").Return(pub, nil)
	w = mockWitnessSignature(w, pub, priv)
	return w
}

func mockWitnessSignature(
	w *mocks.Wallet, pub *btcec.PublicKey, priv *btcec.PrivateKey) *mocks.Wallet {
	call := w.On("WitnessSignature",
		mock.AnythingOfType("*wire.MsgTx"),
		mock.AnythingOfType("int"),
		mock.AnythingOfType("btcutil.Amount"),
		mock.AnythingOfType("[]uint8"),
		pub,
	)

	call.Run(func(args mock.Arguments) {
		tx := args.Get(0).(*wire.MsgTx)
		idx := args.Get(1).(int)
		amt := args.Get(2).(btcutil.Amount)
		sc := args.Get(3).([]uint8)
		sign, err := script.WitnessSignature(tx, idx, int64(amt), sc, priv)
		rargs := make([]interface{}, 2)
		rargs[0] = sign
		rargs[1] = err
		call.ReturnArguments = rargs
	})

	return w
}

func mockWitnessSignatureWithCallback(
	w *mocks.Wallet, pub *btcec.PublicKey, priv *btcec.PrivateKey,
	privkeyConverter wallet.PrivateKeyConverter,
) *mocks.Wallet {
	call := w.On("WitnessSignatureWithCallback",
		mock.AnythingOfType("*wire.MsgTx"),
		mock.AnythingOfType("int"),
		mock.AnythingOfType("btcutil.Amount"),
		mock.AnythingOfType("[]uint8"),
		pub,
		mock.AnythingOfType("wallet.PrivateKeyConverter"),
	)

	call.Run(func(args mock.Arguments) {
		tx := args.Get(0).(*wire.MsgTx)
		idx := args.Get(1).(int)
		amt := args.Get(2).(btcutil.Amount)
		sc := args.Get(3).([]uint8)
		privplus, _ := privkeyConverter(priv)
		sign, err := script.WitnessSignature(tx, idx, int64(amt), sc, privplus)
		rargs := make([]interface{}, 2)
		rargs[0] = sign
		rargs[1] = err
		call.ReturnArguments = rargs
	})

	return w
}

// Hash of block 234439
var testTxID = "14a0810ac680a3eb3f82edc878cea25ec41d6b790744e5daeef"

func mockSelectUnspent(
	w *mocks.Wallet, balance, change btcutil.Amount, err error) *mocks.Wallet {
	utxo := wallet.Utxo{
		TxID:   testTxID,
		Amount: float64(balance) / btcutil.SatoshiPerBitcoin,
	}
	w.On("SelectUnspent",
		mock.Anything, mock.Anything, mock.Anything,
	).Return([]wallet.Utxo{utxo}, change, err)

	return w
}
