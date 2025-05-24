//go:build darwin

package batt

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var MaintainInterval = 60 * time.Second

func maintain(upper, lower int) error {
	slog.Debug("Maintain called", "upper", upper, "lower", lower)

	c := New()
	if err := c.Open(); err != nil {
		return fmt.Errorf("maintain failed to open Apple SMC: %w", err)
	}
	defer c.Close()

	batteryLevel, err := c.GetBatteryLevel()
	if err != nil {
		return fmt.Errorf("maintain failed to get battery level: %w", err)
	}

	isChargingEnabled, err := c.IsChargingEnabled()
	if err != nil {
		return fmt.Errorf("maintain failed to check if charging is enabled: %w", err)
	}

	if batteryLevel >= upper {
		if isChargingEnabled {
			if err := c.DisableCharging(); err != nil {
				return fmt.Errorf("failed to disable charging: %w", err)
			}
		}
		slog.Info(
			"Battery level is at or above upper limit, charging disabled",
			"upper",
			upper,
			"lower",
			lower,
			"level",
			batteryLevel)
	} else if batteryLevel <= lower {
		if !isChargingEnabled {
			if err := c.EnableCharging(); err != nil {
				return fmt.Errorf("failed to enable charging: %w", err)
			}
		}
		slog.Info(
			"Battery level is at or below lower limit, charging enabled",
			"upper",
			upper,
			"lower",
			lower,
			"level",
			batteryLevel)
	} else {
		slog.Info(
			"Battery level is within limits, no action taken",
			"upper",
			upper,
			"lower",
			lower,
			"level",
			batteryLevel)
	}

	return nil
}

func Maintain(upper, lower int) error {
	// Validate upper and lower limits
	if upper < lower {
		return fmt.Errorf("the upper limit %d must be equal or greater than the lower limit %d", upper, lower)
	}

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)
	defer close(stopChan)

	ticker := time.NewTicker(MaintainInterval)
	defer ticker.Stop()

	// Initial maintenance call
	if err := maintain(upper, lower); err != nil {
		return fmt.Errorf("maintain failed: %w", err)
	}

	for {
		select {
		case <-ticker.C:
			if err := maintain(upper, lower); err != nil {
				return fmt.Errorf("maintain failed: %w", err)
			}
		case <-stopChan:
			return nil
		}
	}
}
