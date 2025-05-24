//go:build darwin

package batt

import (
	"log/slog"
)

// IsAdapterEnabled returns whether the adapter is enabled.
func (c *AppleSMC) IsAdapterEnabled() (bool, error) {
	slog.Debug("IsAdapterEnabled called")

	v, err := c.Read(AdapterKey)
	if err != nil {
		return false, err
	}

	ret := len(v.Bytes) == 1 && v.Bytes[0] == AdapterValEnabled
	slog.Debug(
		"IsAdapterEnabled read from SMC succeed",
		slog.String("key", AdapterKey),
		slog.String("val", string(v.Bytes)),
	)

	return ret, nil
}

// EnableAdapter enables the adapter.
func (c *AppleSMC) EnableAdapter() error {
	slog.Debug("EnableAdapter called")

	return c.Write(AdapterKey, []byte{AdapterValEnabled})
}

// DisableAdapter disables the adapter.
func (c *AppleSMC) DisableAdapter() error {
	slog.Debug("DisableAdapter called")

	return c.Write(AdapterKey, []byte{AdapterValDisabled})
}
