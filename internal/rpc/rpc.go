// Package rpc project rpc.go
package rpc

import (
	"github.com/btcsuite/btcd/rpcclient"
)

// NewBtcdRPC returns new rpcclient.Client.
func NewBtcdRPC(url, user, pass string) (*rpcclient.Client, error) {
	// convert url, user, pass strings into connConfig
	connCfg := &rpcclient.ConnConfig{
		Host:         url,
		User:         user,
		Pass:         pass,
		HTTPPostMode: true, // Bitcoin core only supports HTTP POST mode
		DisableTLS:   true, // Bitcoin core does not provide TLS by default
	}

	return rpcclient.New(connCfg, nil)
}
