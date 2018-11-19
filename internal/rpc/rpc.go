// Package rpc project rpc.go
package rpc

import (
	"errors"
	"io/ioutil"
	"net"
	"os"
	"regexp"
	"strings"

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

const (
	defaultHost         = "localhost"
	defaultHTTPPostMode = true
	defaultDisableTLS   = true
)

// NewClient returns Client interface object
func NewClient(cfgPath string) (Client, error) {
	cfg, err := loadConfig(cfgPath)
	if err != nil {
		return nil, err
	}
	return newClient(cfg)
}

func newClient(cfg *rpcclient.ConnConfig) (*rpcclient.Client, error) {
	return rpcclient.New(cfg, nil)
}

func loadConfig(cfgPath string) (*rpcclient.ConnConfig, error) {
	cfgFile, err := os.Open(cfgPath)
	if err != nil {
		return nil, err
	}
	defer cfgFile.Close()

	content, err := ioutil.ReadAll(cfgFile)
	if err != nil {
		return nil, err
	}

	// Extract the rpcuser
	rpcUserRegexp, err := regexp.Compile(`(?m)^\s*rpcuser=([^\s]+)`)
	if err != nil {
		return nil, err
	}
	userSubmatches := rpcUserRegexp.FindSubmatch(content)
	if userSubmatches == nil {
		return nil, errors.New("rpcuser isn't set in config file")
	}
	user := strings.Split(string(userSubmatches[0]), "=")[1]

	// Extract the rpcpassword
	rpcPassRegexp, err := regexp.Compile(`(?m)^\s*rpcpassword=([^\s]+)`)
	if err != nil {
		return nil, err
	}
	passSubmatches := rpcPassRegexp.FindSubmatch(content)
	if passSubmatches == nil {
		return nil, errors.New("rpcpassword isn't set in config file")
	}
	pass := strings.Split(string(passSubmatches[0]), "=")[1]

	// Extract the regtest
	regtestRegexp, err := regexp.Compile(`(?m)^\s*regtest=([^\s]+)`)
	if err != nil {
		return nil, err
	}
	regtestSubmatches := regtestRegexp.FindSubmatch(content)
	regtest := strings.Split(string(regtestSubmatches[0]), "=")[1]
	var useRegTest bool
	if regtestSubmatches != nil && regtest == "1" {
		useRegTest = true
	}

	cfg := &rpcclient.ConnConfig{
		Host:         appendPort(defaultHost, useRegTest),
		User:         user,
		Pass:         pass,
		HTTPPostMode: defaultHTTPPostMode, // Bitcoin core only supports HTTP POST mode
		DisableTLS:   defaultDisableTLS,   // Bitcoin core does not provide TLS by default
	}
	return cfg, nil
}

func appendPort(addr string, useRegTest bool) string {
	port := ""
	if useRegTest {
		port = "18443"
	}

	return net.JoinHostPort(addr, port)
}
