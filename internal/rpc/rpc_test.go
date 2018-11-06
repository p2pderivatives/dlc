package rpc

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testport = "localhost:18332" //18443 for regnet, 18332 for testnet3
	testuser = "akek"
	testpass = "akek"
)

func TestNewBtcdRPC(t *testing.T) {
	client, err := NewBtcdRPC(testport, testuser, testpass)
	defer client.Shutdown()
	if err != nil {
		log.Fatal(err)
	}
	assert.NotNil(t, client)

	// Get the current block count.
	blockCount, err := client.GetBlockCount()
	if err != nil {
		log.Fatal(err)
	}
	info2, _ := client.GetBlockChainInfo()
	info3, _ := client.ListUnspent()
	info4, _ := client.ListAccounts()

	log.Printf("Block count: %d \n\n", blockCount)
	log.Printf("BLOCKCHAIN INFO: %+v \n\n", info2)
	log.Printf("LIST UNSPENT: %+v \n\n", info3)
	log.Printf("LIST accounts: %+v \n\n", info4)

}
