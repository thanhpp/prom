package timerpc

import (
	"fmt"
	"time"
)

func (t *Time) Scan(value interface{}) error {
	ti, ok := value.(time.Time)
	if !ok {
		return fmt.Errorf("Failed to cast time. Value: %v", value)
	}

	return nil
}
