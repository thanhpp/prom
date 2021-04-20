package timerpc

import (
	"database/sql/driver"
	"fmt"
	"time"
)

func (t *Time) Scan(value interface{}) error {
	ti, ok := value.(time.Time)
	if !ok {
		return fmt.Errorf("Failed to cast time. Value: %v", value)
	}
	t.Seconds = int64(ti.Unix())
	t.Nanos = int32(ti.Nanosecond())

	return nil
}

func (t *Time) Value() (driver.Value, error) {
	if t == nil {
		return nil, nil
	}
	return time.Unix(t.Seconds, int64(t.Nanos)), nil
}
