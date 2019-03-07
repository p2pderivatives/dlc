package wallet

import (
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/p2pderivatives/dlc/internal/rpc"
)

// Wallet is an interface that provides access to manage pubkey addresses and
// sign scripts of managed addressesc using private key. It also manags utxos.
type Wallet interface {
	NewPubkey() (*btcec.PublicKey, error)

	// NewAddress creates a new address
	NewAddress() (btcutil.Address, error)

	// WitnessSignature returns witness signature for a given txin and pubkey
	WitnessSignature(
		tx *wire.MsgTx, idx int, amt btcutil.Amount, sc []byte, pub *btcec.PublicKey,
	) (sign []byte, err error)

	// WitnessSignatureWithCallback does the same with WitnessSignature do
	// applying a given func to private key before calculating signature
	WitnessSignatureWithCallback(
		tx *wire.MsgTx, idx int, amt btcutil.Amount, sc []byte, pub *btcec.PublicKey,
		privkeyConverter PrivateKeyConverter,
	) (sign []byte, err error)

	// WitnessSignTxByIdxs returns witness signatures for txins specified by idxs
	WitnessSignTxByIdxs(tx *wire.MsgTx, idxs []int) ([]wire.TxWitness, error)

	// SelectUtxos selects utxos for requested amount
	// by considering additional fee per txin and txout
	SelectUnspent(
		amt, feePerTxIn, feePerTxOut btcutil.Amount,
	) (utxos []Utxo, change btcutil.Amount, err error)
	// Unlock unlocks address manager
	Unlock(privPass []byte) error

	// TODO: remove this interface after fixing wallet.Open
	// SetRPCClient sets rpcclient
	SetRPCClient(rpc.Client)

	// methods delegating to RPC Client
	ListUnspent() (utxos []Utxo, err error)
	SendRawTransaction(tx *wire.MsgTx) (*chainhash.Hash, error)

	Close() error
}

// Utxo is an unspent transaction output
type Utxo = btcjson.ListUnspentResult

// PrivateKeyConverter is a callback func applied to private key before creating witness signature
type PrivateKeyConverter func(*btcec.PrivateKey) (*btcec.PrivateKey, error)
