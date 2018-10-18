package oracle

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcutil/hdkeychain"
)

// TimeFormat is a format of settlement time
const TimeFormat = "20060102"

// Oracle is a struct
type Oracle struct {
	name   string                  // display name
	digit  int                     // digit
	extKey *hdkeychain.ExtendedKey // extended key
}

// New creates a oracle
func New(name string, params chaincfg.Params, digit int) (*Oracle, error) {
	extKey, err := randomExtKey(name, params)
	if err != nil {
		return nil, err
	}

	oracle := &Oracle{name: name, digit: digit, extKey: extKey}
	return oracle, nil
}

// randomExtKey creates oracle's random master key
func randomExtKey(name string, params chaincfg.Params) (*hdkeychain.ExtendedKey, error) {
	// TODO: add random logic
	seed := chainhash.DoubleHashB([]byte(name))
	return hdkeychain.NewMaster(seed, &params)
}
