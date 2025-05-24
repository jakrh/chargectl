//go:build darwin

package gosmc

import (
	"bytes"
	"os"
	"testing"
)

func TestSMC(t *testing.T) {
	c := New()
	err := c.Open()
	if err != nil {
		t.Fatal("Failed to open SMC:", err)
	}
	defer c.Close()

	// **requires root access** enable charging
	if os.Getuid() != 0 {
		t.Skip("Skipping test that requires root access")
	}

	want := struct {
		key      string
		val      []byte
		dataType string
	}{
		key:      "CH0C",
		val:      []byte{0x0},
		dataType: "hex_",
	}

	err = c.Write(want.key, want.val)
	if err != nil {
		t.Errorf("Failed to write to SMC: %v", err)
	}

	v, err := c.Read(want.key)
	if err != nil {
		t.Errorf("Failed to read from SMC by key: %s, error: %v", want.key, err)
	}

	if v.Key != want.key {
		t.Errorf("Expected key to not be %s, got %s", want.key, v.Key)
	}

	if v.DataType != want.dataType {
		t.Errorf("Expected data type to be %s, got %s", want.dataType, v.DataType)
	}

	if !bytes.Equal(v.Bytes, want.val) {
		t.Errorf("Expected bytes to be %v, got %v", want.val, v.Bytes)
	}
}
