package rpc

import (
	"log"
	"testing"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/dgarage/dlc/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

func mockListUnspent(w *mocks.Wallet, err error) *mocks.Wallet {
	utxo := btcjson.ListUnspentResult{
		TxID:          "ce9d930c2664547ad8aba6944c8047321bde0c1c1d6551c41ebb8d9ad975dd0b",
		Vout:          uint32(0),
		Address:       "tb1qds49lkplvws9q4df04e5e9nq5d6asnkkhna8hg",
		Account:       "",
		ScriptPubKey:  "00146c2a5fd83f63a05055a97d734c9660a375d84ed6",
		RedeemScript:  "",
		Amount:        float64(0.31864472),
		Confirmations: int64(30006),
		Spendable:     true,
	}

	w.On("ListUnspent").Return([]btcjson.ListUnspentResult{utxo}, err)

	return w
}

func mockImportAddress(w *mocks.Wallet, address string, err error) *mocks.Wallet {
	w.On("ImportAddress",
		mock.AnythingOfType("string"),
	).Return(err)

	return w
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
