package oracle

import "github.com/btcsuite/btcd/btcec"

// PubkeySet contains oracle's pub key and keys for all rate
type PubkeySet struct {
	Pubkey           *btcec.PublicKey
	CommittedRpoints []*btcec.PublicKey
}

// SignSet contains fixed messages and signs
type SignSet struct {
	Msgs  [][]byte
	Signs [][]byte
}
