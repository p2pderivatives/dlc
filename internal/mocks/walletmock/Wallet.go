// Code generated by mockery v1.0.0. DO NOT EDIT.

package walletmock

import btcec "github.com/btcsuite/btcd/btcec"
import btcjson "github.com/btcsuite/btcd/btcjson"
import btcutil "github.com/btcsuite/btcutil"
import chainhash "github.com/btcsuite/btcd/chaincfg/chainhash"
import mock "github.com/stretchr/testify/mock"
import rpc "github.com/p2pderivatives/dlc/internal/rpc"
import wallet "github.com/p2pderivatives/dlc/pkg/wallet"
import wire "github.com/btcsuite/btcd/wire"

// Wallet is an autogenerated mock type for the Wallet type
type Wallet struct {
	mock.Mock
}

// Close provides a mock function with given fields:
func (_m *Wallet) Close() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ListUnspent provides a mock function with given fields:
func (_m *Wallet) ListUnspent() ([]btcjson.ListUnspentResult, error) {
	ret := _m.Called()

	var r0 []btcjson.ListUnspentResult
	if rf, ok := ret.Get(0).(func() []btcjson.ListUnspentResult); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]btcjson.ListUnspentResult)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewAddress provides a mock function with given fields:
func (_m *Wallet) NewAddress() (btcutil.Address, error) {
	ret := _m.Called()

	var r0 btcutil.Address
	if rf, ok := ret.Get(0).(func() btcutil.Address); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(btcutil.Address)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewPubkey provides a mock function with given fields:
func (_m *Wallet) NewPubkey() (*btcec.PublicKey, error) {
	ret := _m.Called()

	var r0 *btcec.PublicKey
	if rf, ok := ret.Get(0).(func() *btcec.PublicKey); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*btcec.PublicKey)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SelectUnspent provides a mock function with given fields: amt, feePerTxIn, feePerTxOut
func (_m *Wallet) SelectUnspent(amt btcutil.Amount, feePerTxIn btcutil.Amount, feePerTxOut btcutil.Amount) ([]btcjson.ListUnspentResult, btcutil.Amount, error) {
	ret := _m.Called(amt, feePerTxIn, feePerTxOut)

	var r0 []btcjson.ListUnspentResult
	if rf, ok := ret.Get(0).(func(btcutil.Amount, btcutil.Amount, btcutil.Amount) []btcjson.ListUnspentResult); ok {
		r0 = rf(amt, feePerTxIn, feePerTxOut)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]btcjson.ListUnspentResult)
		}
	}

	var r1 btcutil.Amount
	if rf, ok := ret.Get(1).(func(btcutil.Amount, btcutil.Amount, btcutil.Amount) btcutil.Amount); ok {
		r1 = rf(amt, feePerTxIn, feePerTxOut)
	} else {
		r1 = ret.Get(1).(btcutil.Amount)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(btcutil.Amount, btcutil.Amount, btcutil.Amount) error); ok {
		r2 = rf(amt, feePerTxIn, feePerTxOut)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// SendRawTransaction provides a mock function with given fields: tx
func (_m *Wallet) SendRawTransaction(tx *wire.MsgTx) (*chainhash.Hash, error) {
	ret := _m.Called(tx)

	var r0 *chainhash.Hash
	if rf, ok := ret.Get(0).(func(*wire.MsgTx) *chainhash.Hash); ok {
		r0 = rf(tx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*chainhash.Hash)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*wire.MsgTx) error); ok {
		r1 = rf(tx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetRPCClient provides a mock function with given fields: _a0
func (_m *Wallet) SetRPCClient(_a0 rpc.Client) {
	_m.Called(_a0)
}

// Unlock provides a mock function with given fields: privPass
func (_m *Wallet) Unlock(privPass []byte) error {
	ret := _m.Called(privPass)

	var r0 error
	if rf, ok := ret.Get(0).(func([]byte) error); ok {
		r0 = rf(privPass)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// WitnessSignTxByIdxs provides a mock function with given fields: tx, idxs
func (_m *Wallet) WitnessSignTxByIdxs(tx *wire.MsgTx, idxs []int) ([]wire.TxWitness, error) {
	ret := _m.Called(tx, idxs)

	var r0 []wire.TxWitness
	if rf, ok := ret.Get(0).(func(*wire.MsgTx, []int) []wire.TxWitness); ok {
		r0 = rf(tx, idxs)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]wire.TxWitness)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*wire.MsgTx, []int) error); ok {
		r1 = rf(tx, idxs)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// WitnessSignature provides a mock function with given fields: tx, idx, amt, sc, pub
func (_m *Wallet) WitnessSignature(tx *wire.MsgTx, idx int, amt btcutil.Amount, sc []byte, pub *btcec.PublicKey) ([]byte, error) {
	ret := _m.Called(tx, idx, amt, sc, pub)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(*wire.MsgTx, int, btcutil.Amount, []byte, *btcec.PublicKey) []byte); ok {
		r0 = rf(tx, idx, amt, sc, pub)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*wire.MsgTx, int, btcutil.Amount, []byte, *btcec.PublicKey) error); ok {
		r1 = rf(tx, idx, amt, sc, pub)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// WitnessSignatureWithCallback provides a mock function with given fields: tx, idx, amt, sc, pub, privkeyConverter
func (_m *Wallet) WitnessSignatureWithCallback(tx *wire.MsgTx, idx int, amt btcutil.Amount, sc []byte, pub *btcec.PublicKey, privkeyConverter wallet.PrivateKeyConverter) ([]byte, error) {
	ret := _m.Called(tx, idx, amt, sc, pub, privkeyConverter)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(*wire.MsgTx, int, btcutil.Amount, []byte, *btcec.PublicKey, wallet.PrivateKeyConverter) []byte); ok {
		r0 = rf(tx, idx, amt, sc, pub, privkeyConverter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*wire.MsgTx, int, btcutil.Amount, []byte, *btcec.PublicKey, wallet.PrivateKeyConverter) error); ok {
		r1 = rf(tx, idx, amt, sc, pub, privkeyConverter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
