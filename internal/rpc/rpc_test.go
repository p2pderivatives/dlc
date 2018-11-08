package rpc

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testport = "localhost:18433" //18443 for regnet, 18332 for testnet3
	testuser = "username"
	testpass = "password"
)

func TestNewClient(t *testing.T) {
	client, err := NewClient(testport, testuser, testpass)
	//defer client.Shutdown()
	if err != nil {
		log.Fatal(err)
	}
	assert.NotNil(t, client)
}

// func setupRegNet() {
// 	connCfg := &rpcclient.ConnConfig{
// 		Host:         testport,
// 		User:         testuser,
// 		Pass:         testpass,
// 		HTTPPostMode: true, // Bitcoin core only supports HTTP POST mode
// 		DisableTLS:   true, // Bitcoin core does not provide TLS by default
// 	}
// 	c, _ := rpcclient.New(connCfg, nil)
// 	c.Generate(150)
// 	blockCount, _ := c.GetBlockCount()
// 	info2, _ := c.GetBlockChainInfo()
// 	info3, _ := c.ListUnspent()

// 	log.Printf("Block count: %d \n\n", blockCount)
// 	log.Printf("BLOCKCHAIN INFO: %+v \n\n", info2)
// 	log.Printf("LIST UNSPENT: %+v \n\n", info3)
// }
