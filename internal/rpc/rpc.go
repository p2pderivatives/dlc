// Package rpc project rpc.go
package rpc

import (
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
)

// Client is an interface that provides access to certain methods of type rpcclient.Client
type Client interface {
	ListUnspent() ([]btcjson.ListUnspentResult, error)
	ImportAddress(address string) error
	SendRawTransaction(tx *wire.MsgTx, allowHighFees bool) (*chainhash.Hash, error)
	// TODO: add Shutdown func
}

// NewClient returns Client interface object
func NewClient(url, user, pass string) (Client, error) {
	return newClient(url, user, pass)
}

func newClient(url, user, pass string) (*rpcclient.Client, error) {
	connCfg := &rpcclient.ConnConfig{
		Host:         url,
		User:         user,
		Pass:         pass,
		HTTPPostMode: true, // Bitcoin core only supports HTTP POST mode
		DisableTLS:   true, // Bitcoin core does not provide TLS by default
	}

	return rpcclient.New(connCfg, nil)
}
