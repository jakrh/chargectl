//go:build darwin

package batt

// Various SMC keys for arm64 (Apple Silicon)
const (
	ChargingKey1    = "CH0B"
	ChargingKey2    = "CH0C"
	AdapterKey      = "CH0I"
	BatteryLevelKey = "BUIC"

	AdapterValEnabled   byte = 0x0
	AdapterValDisabled  byte = 0x1
	ChargingValEnabled  byte = 0x0
	ChargingValDisabled byte = 0x2
)

var ChargingKeys = []string{ChargingKey1, ChargingKey2}
