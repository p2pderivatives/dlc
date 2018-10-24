// Package rpc project rpc.go
package rpc

import (
	"log"

	"github.com/btcsuite/btcd/rpcclient"
)

// BtcRPC is request info.
type BtcRPC struct {
	URL  string // bitcoin full node endpoint url
	User string // rpcuser
	Pass string // rpcpassword
	View bool   // If true, the log is displayed.
}

// BtcRPCRequest is request parameters.
type BtcRPCRequest struct {
	// bitcoin rpc request format
	Jsonrpc string        `json:"jsonrpc,"`
	ID      string        `json:"id,"`
	Method  string        `json:"method,"`
	Params  []interface{} `json:"params,"`
}

// Response is response details.
type Response struct {
	Result interface{} `json:"result,"`
	Error  interface{} `json:"error,"`
	ID     string      `json:"id,"`
}

// Error is error details.
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

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

func (rpc *BtcRPC) log(format string, v ...interface{}) {
	if rpc.View {
		log.Printf(format, v...)
	}
}
