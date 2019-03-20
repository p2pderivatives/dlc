package dlc

import (
	"errors"
	"fmt"

	"github.com/btcsuite/btcutil"
)

// NotEnoughFeesError is an error that txins aren't enough for txouts
type NotEnoughFeesError struct {
	error
}

func newNotEnoughFeesError(in, fee btcutil.Amount) *NotEnoughFeesError {
	msg := fmt.Sprintf("TxFee isn't enough. txins: %d, fee: %d", in, fee)
	return &NotEnoughFeesError{error: errors.New(msg)}
}

// NoFixedDealError is an error for a case when no deals has been fixed yet
type NoFixedDealError struct {
	error
}

func newNoFixedDealError() *NoFixedDealError {
	msg := "No deal has been fixed"
	return &NoFixedDealError{error: errors.New(msg)}
}
