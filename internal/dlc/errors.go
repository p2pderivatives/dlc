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

// CETTakeNothingError is an error for invalid CET
// Doesn't make sense to create a CTE that takes nothing
type CETTakeNothingError struct {
	error
}

func newCETTakeNothingError(msg string) *CETTakeNothingError {
	return &CETTakeNothingError{error: errors.New(msg)}
}

// NoFixedDealError is an error for a case when no deals has been fixed yet
type NoFixedDealError struct {
	error
}

func newNoFixedDealError() *NoFixedDealError {
	msg := "No deal has been fixed"
	return &NoFixedDealError{error: errors.New(msg)}
}
