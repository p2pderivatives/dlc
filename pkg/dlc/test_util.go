package dlc

import (
	"time"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/p2pderivatives/dlc/internal/mocks/walletmock"
	"github.com/p2pderivatives/dlc/internal/test"
	"github.com/p2pderivatives/dlc/pkg/script"
	"github.com/p2pderivatives/dlc/pkg/wallet"
	"github.com/stretchr/testify/mock"
)

// setup mocke wallet
func setupTestWallet() *walletmock.Wallet {
	w := &walletmock.Wallet{}
	w = mockNewAddress(w)

	// pubkey for fund script
	priv, pub := test.RandKeys()
	w.On("NewPubkey").Return(pub, nil)
	w = mockWitnessSignature(w, pub, priv)

	return w
}

func mockNewAddress(w *walletmock.Wallet) *walletmock.Wallet {
	_, pub := test.RandKeys()
	pubKeyHash := btcutil.Hash160(pub.SerializeCompressed())
	addr, _ := btcutil.NewAddressWitnessPubKeyHash(
		pubKeyHash, &chaincfg.RegressionNetParams)
	w.On("NewAddress").Return(addr, nil)
	return w
}

func mockWitnessSignature(
	w *walletmock.Wallet, pub *btcec.PublicKey, priv *btcec.PrivateKey) *walletmock.Wallet {
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
	w *walletmock.Wallet, pub *btcec.PublicKey, priv *btcec.PrivateKey,
	privkeyConverter wallet.PrivateKeyConverter,
) *walletmock.Wallet {
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
	w *walletmock.Wallet, balance, change btcutil.Amount, err error) *walletmock.Wallet {
	utxo := wallet.Utxo{
		TxID:   testTxID,
		Amount: float64(balance) / btcutil.SatoshiPerBitcoin,
	}
	w.On("SelectUnspent",
		mock.Anything, mock.Anything, mock.Anything,
	).Return([]wallet.Utxo{utxo}, change, err)

	return w
}

func newTestConditions() *Conditions {
	conds, _ := NewConditions(time.Now(), 1, 1, 1, 1, 1, []*Deal{})
	return conds
}

// stepSendRequirments send pubkey, utoxs, change address
func stepSendRequirments(b1, b2 *Builder) (*Builder, *Builder) {
	// b1 -> b2
	p1, _ := b1.PublicKey()
	u1 := b1.Utxos()
	caddr1 := b1.ChangeAddress()
	b2.AcceptPubkey(p1)
	b2.AcceptUtxos(u1)
	b2.AcceptsChangeAdderss(caddr1)

	return b1, b2
}
