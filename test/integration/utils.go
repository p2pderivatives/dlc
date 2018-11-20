package integration

import "github.com/btcsuite/btcutil"

// Faucet sends bitcion to a given address
func Faucet(addr btcutil.Address, amt btcutil.Amount) error {
	c, err := NewRPCClient()
	if err != nil {
		return err
	}
	_, err = c.SendToAddress(addr, amt)
	if err != nil {
		return err
	}

	_, err = c.Generate(1)
	return err
}
