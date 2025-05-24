//go:build darwin

package batt

import (
	"fmt"
	"log/slog"
)

// GetBatteryLevel returns the battery level.
func (c *AppleSMC) GetBatteryLevel() (int, error) {
	slog.Debug("GetBatteryLevel called")

	v, err := c.Read(BatteryLevelKey)
	if err != nil {
		return 0, err
	}

	if len(v.Bytes) != 1 {
		return 0, fmt.Errorf("incorrect data length %d!=1", len(v.Bytes))
	}

	return int(v.Bytes[0]), nil
}
