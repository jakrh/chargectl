//go:build darwin

package batt

// Various SMC keys for amd64 (Intel x64).
const (
	ChargingKey     = "BCLM"
	AdapterKey      = "CH0K"
	BatteryLevelKey = "BBIF"

	AdapterValEnabled   byte = 0x0
	AdapterValDisabled  byte = 0x1
	ChargingValEnabled  byte = 0x64
	ChargingValDisabled byte = 0x0a
)

var ChargingKeys = []string{ChargingKey}
