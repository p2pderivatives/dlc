package dlc

import "errors"

// CETTakeNothingError is an error for invalid CET
// Doesn't make sense to create a CTE that takes nothing
type CETTakeNothingError struct {
	error
}

func newCETTakeNothingError(msg string) *CETTakeNothingError {
	return &CETTakeNothingError{error: errors.New(msg)}
}
