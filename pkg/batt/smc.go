//go:build darwin

package batt

import (
	"log/slog"

	"chargectl/pkg/gosmc"
)

// AppleSMC is a wrapper of gosmc.Connection.
type AppleSMC struct {
	conn gosmc.Conn
}

// New returns a new AppleSMC.
func New() *AppleSMC {
	return &AppleSMC{
		conn: gosmc.New(),
	}
}

// NewMock returns a new mocked AppleSMC with prefill values.
func NewMock(prefillValues map[string][]byte) *AppleSMC {
	conn := gosmc.NewMockConn()

	for key, value := range prefillValues {
		err := conn.Write(key, value)
		if err != nil {
			panic(err)
		}
	}

	return &AppleSMC{
		conn: conn,
	}
}

// Open opens the connection.
func (c *AppleSMC) Open() error {
	return c.conn.Open()
}

// Close closes the connection.
func (c *AppleSMC) Close() error {
	return c.conn.Close()
}

// Read reads a value from SMC.
func (c *AppleSMC) Read(key string) (gosmc.SMCVal, error) {
	slog.Debug("Trying to read from SMC", slog.String("key", key))

	v, err := c.conn.Read(key)
	if err != nil {
		return v, err
	}

	slog.Debug("Read from SMC succeed", slog.String("key", key), slog.String("val", string(v.Bytes)))

	return v, nil
}

// Write writes a value to SMC.
func (c *AppleSMC) Write(key string, value []byte) error {
	slog.Debug("Trying to write to SMC", slog.String("key", key), slog.String("val", string(value)))

	err := c.conn.Write(key, value)
	if err != nil {
		return err
	}

	slog.Debug("Write to SMC succeed", slog.String("key", key), slog.String("val", string(value)))

	return nil
}
