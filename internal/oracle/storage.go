package oracle

import (
	"fmt"
	"time"
)

// TimeFormat is a format of settlement time
const TimeFormat = "2006-01-02 15:04:05"

type memdb struct {
	values map[string][]int
}

func (oracle *Oracle) valuesAt(ftime time.Time) ([]int, error) {
	key := ftime.Format(TimeFormat)
	vals, ok := oracle.db.values[key]
	if !ok {
		return []int{}, fmt.Errorf("not found values at %s", key)
	}
	return vals, nil
}

func (oracle *Oracle) fixValues(ftime time.Time, values []int) error {
	size := oracle.nRpoints
	if len(values) != size {
		return fmt.Errorf("invalid values size. expected %d, but got %d", size, len(values))
	}
	key := ftime.Format(TimeFormat)
	oracle.db.values[key] = values
	return nil
}
