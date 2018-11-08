// Package rpc project rpc.go
package rpc

import (
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
)

type Client interface {
	ListUnspent() ([]btcjson.ListUnspentResult, error)
	//ImportManagedAddress(waddrmgr.ManagedAddress) error
	ImportAddress(address string) error
	SendRawTransaction(tx *wire.MsgTx, allowHighFees bool) (*chainhash.Hash, error)
}

// type client struct {
// 	client *rpcclient.Client
// }

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

	// client, err := rpcclient.New(connCfg, nil)
	// if err != nil {
	// 	return nil, err
	// }

	// c := &client{
	// 	client: client,
	// }

	return rpcclient.New(connCfg, nil)
	// return c, err
}

// ImportAddress get the public address from the passed ManagedAddress and imports it
// func (c *client) ImportManagedAddress(maddr waddrmgr.ManagedAddress) error {
// 	return c.client.ImportAddress(maddr.(waddrmgr.ManagedPubKeyAddress).ExportPubKey())
// }
