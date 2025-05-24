//go:build darwin

package batt

import (
	"fmt"
	"log/slog"
)

// IsChargingEnabled returns whether charging is enabled.
func (c *AppleSMC) IsChargingEnabled() (bool, error) {
	slog.Debug("IsChargingEnabled called")

	ret := true
	for _, key := range ChargingKeys {
		v, err := c.Read(key)
		if err != nil {
			return false, fmt.Errorf("failed to check charging status for key %s: %w", key, err)
		}

		ret = ret && len(v.Bytes) == 1 && v.Bytes[0] == ChargingValEnabled
		slog.Debug(
			"IsChargingEnabled read from SMC succeed",
			slog.String("key", key),
			slog.String("val", string(v.Bytes)))
	}

	return ret, nil
}

// EnableCharging enables charging.
func (c *AppleSMC) EnableCharging() error {
	slog.Debug("EnableCharging called")

	for _, key := range ChargingKeys {
		err := c.Write(key, []byte{ChargingValEnabled})
		if err != nil {
			return fmt.Errorf("failed to enable charging for key %s: %w", key, err)
		}
	}

	return c.EnableAdapter()
}

// DisableCharging disables charging.
func (c *AppleSMC) DisableCharging() error {
	slog.Debug("DisableCharging called")

	for _, key := range ChargingKeys {
		err := c.Write(key, []byte{ChargingValDisabled})
		if err != nil {
			return fmt.Errorf("failed to disable charging for key %s: %w", key, err)
		}
	}

	return nil
}
