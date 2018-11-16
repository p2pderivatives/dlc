package integration

import (
	"io/ioutil"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/dgarage/dlc/internal/dlc"
	"github.com/dgarage/dlc/internal/wallet"
)

// Contractor is a contractor
type Contractor struct {
	Name       string
	Wallet     wallet.Wallet
	DLCBuilder *dlc.Builder
	pubpass    []byte
	privpass   []byte
}

// newContractor creates a contractor for integration tests
func newContractor(name string) (*Contractor, error) {
	c := &Contractor{Name: name}
	err := c.createWallet()
	return c, err
}

// createWallet creates a wallet
func (c *Contractor) createWallet() error {
	params := &chaincfg.RegressionNetParams
	c.pubpass = []byte("pubpass")
	c.privpass = []byte("privpass")

	// generate random seed
	seed, err := hdkeychain.GenerateSeed(
		hdkeychain.RecommendedSeedLen)
	if err != nil {
		return err
	}

	// create wallet dbdir
	walletDir, err := ioutil.TempDir("", "dlcwallet")
	if err != nil {
		return err
	}

	// create wallet
	walletName := c.Name
	c.Wallet, err = wallet.CreateWallet(
		params, seed, c.pubpass, c.privpass, walletDir, walletName)
	return err
}

func (c *Contractor) createDLCBuilder(
	conds *dlc.Conditions, p dlc.Contractor) {
	c.DLCBuilder = dlc.NewBuilder(p, c.Wallet, conds)
}
