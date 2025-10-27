package netinterface

import (
	"bytes"
	"strings"
	"testing"
)

func TestReadWrite(t *testing.T) {
	input := `auto eth0
iface eth0 inet dhcp
`

	file, err := Read(strings.NewReader(input))
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	if len(file.Interfaces) != 1 {
		t.Fatalf("Expected 1 interface, got %d", len(file.Interfaces))
	}

	var buf bytes.Buffer
	if err := Write(&buf, file); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "auto eth0") {
		t.Error("Output missing 'auto eth0'")
	}
	if !strings.Contains(output, "iface eth0 inet dhcp") {
		t.Error("Output missing 'iface eth0 inet dhcp'")
	}
}

func TestReadFile(t *testing.T) {
	file, err := ReadFile("testdata/interfaces")
	if err == nil {
		if len(file.Interfaces) <= 0 {
			t.Error("Expected interfaces to be parsed")
		}
	}
}

func TestWriteFile(t *testing.T) {
	file := &File{
		Interfaces: []Interface{
			{
				Name:   "eth0",
				Family: "inet",
				Method: "dhcp",
			},
		},
	}

	err := WriteFile("/tmp/test_interfaces", file)
	if err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}
}

func TestDefaultInterfacesPath(t *testing.T) {
	if DefaultInterfacesPath != "/etc/network/interfaces" {
		t.Errorf("Expected DefaultInterfacesPath to be '/etc/network/interfaces', got '%s'", DefaultInterfacesPath)
	}
}
