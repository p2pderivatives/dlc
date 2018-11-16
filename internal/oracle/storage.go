package oracle

import (
	"fmt"
	"time"
)

// TimeFormat is a format of settlement time
const TimeFormat = "2006-01-02 15:04:05"

type memdb struct {
	msgs map[string][][]byte
}

// InitDB initialized oracle's DB
func (o *Oracle) InitDB() {
	msgs := make(map[string][][]byte)
	o.db = &memdb{msgs: msgs}
}

func (o *Oracle) dbReady() bool {
	return (o.db != nil) && (o.db.msgs != nil)
}

func (o *Oracle) msgsAt(ftime time.Time) ([][]byte, error) {
	if !o.dbReady() {
		return [][]byte{}, fmt.Errorf("DB isn't ready")
	}

	key := ftime.Format(TimeFormat)
	vals, ok := o.db.msgs[key]
	if !ok {
		return [][]byte{}, fmt.Errorf("not found messages at %s", key)
	}
	return vals, nil
}

// FixMsgs fixes messsages at a specified time
func (o *Oracle) FixMsgs(ftime time.Time, msgs [][]byte) error {
	if !o.dbReady() {
		return fmt.Errorf("DB isn't ready")
	}

	size := o.nRpoints
	if len(msgs) != size {
		return fmt.Errorf("invalid messages size. expected %d, but got %d", size, len(msgs))
	}
	key := ftime.Format(TimeFormat)
	o.db.msgs[key] = msgs

	return nil
}
