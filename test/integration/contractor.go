package integration

import (
	"github.com/btcsuite/btcutil"
	"github.com/p2pderivatives/dlc/pkg/dlc"
	"github.com/p2pderivatives/dlc/pkg/wallet"
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
func (c *Contractor) createWallet() (err error) {
	c.pubpass = []byte("pubpass")
	c.privpass = []byte("privpass")
	c.Wallet, err = newWallet(c.Name, c.pubpass, c.privpass)
	return
}

func (c *Contractor) createDLCBuilder(
	conds *dlc.Conditions, p dlc.Contractor) {
	c.DLCBuilder = dlc.NewBuilder(p, c.Wallet, conds)
}

func (c *Contractor) unlockWallet() {
	c.Wallet.Unlock(c.privpass)
}

func (c *Contractor) balance() (total btcutil.Amount, err error) {
	utxos, err := c.Wallet.ListUnspent()
	if err != nil {
		return
	}

	for _, utxo := range utxos {
		amt, err := btcutil.NewAmount(utxo.Amount)
		if err != nil {
			return total, err
		}
		total += amt
	}

	return total, nil
}
