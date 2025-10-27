package netinterface

import (
	"bytes"
	"net"
	"strings"
	"testing"
)

func TestGenerate(t *testing.T) {
	file := &File{
		Sources: []string{"/etc/network/interfaces.d/*"},
		Interfaces: []Interface{
			{
				Name:     "eth0",
				Auto:     true,
				Hotplug:  true,
				Family:   "inet",
				Method:   "static",
				Address:  "192.168.1.10",
				Netmask:  "255.255.255.0",
				Gateway:  "192.168.1.1",
				HwAddress: "ethernet 00:11:22:33:44:55",
				DNS:      []net.IP{net.ParseIP("8.8.8.8"), net.ParseIP("8.8.4.4")},
				DNSSearch: []string{"example.com"},
				Metric:   intPtr(100),
				MTU:      intPtr(1500),
			},
			{
				Name:   "lo",
				Auto:   true,
				Family: "inet",
				Method: "loopback",
			},
		},
	}

	var buf bytes.Buffer
	if err := GenerateTo(&buf, file); err != nil {
		t.Fatalf("GenerateTo failed: %v", err)
	}

	output := buf.String()

	if !strings.Contains(output, "source /etc/network/interfaces.d/*") {
		t.Error("Output should contain source directive")
	}

	if !strings.Contains(output, "auto eth0") {
		t.Error("Output should contain 'auto eth0'")
	}

	if !strings.Contains(output, "allow-hotplug eth0") {
		t.Error("Output should contain 'allow-hotplug eth0'")
	}

	if !strings.Contains(output, "iface eth0 inet static") {
		t.Error("Output should contain 'iface eth0 inet static'")
	}

	if !strings.Contains(output, "address 192.168.1.10") {
		t.Error("Output should contain address")
	}

	if !strings.Contains(output, "hwaddress ethernet 00:11:22:33:44:55") {
		t.Error("Output should contain hwaddress")
	}

	if !strings.Contains(output, "dns-nameservers 8.8.8.8 8.8.4.4") {
		t.Error("Output should contain dns-nameservers")
	}

	if !strings.Contains(output, "auto lo") {
		t.Error("Output should contain 'auto lo'")
	}

	if !strings.Contains(output, "iface lo inet loopback") {
		t.Error("Output should contain 'iface lo inet loopback'")
	}
}

func TestGenerateCustom(t *testing.T) {
	file := &File{
		Interfaces: []Interface{
			{
				Name:   "eth0",
				Family: "inet",
				Method: "static",
				Custom: map[string][]string{
					"custom-opt":  []string{"value1", "value2"},
					"another-opt": []string{"test"},
				},
			},
		},
	}

	var buf bytes.Buffer
	if err := GenerateTo(&buf, file); err != nil {
		t.Fatalf("GenerateTo failed: %v", err)
	}

	output := buf.String()

	if !strings.Contains(output, "custom-opt value1") {
		t.Error("Output should contain custom-opt value1")
	}

	if !strings.Contains(output, "custom-opt value2") {
		t.Error("Output should contain custom-opt value2")
	}

	if !strings.Contains(output, "another-opt test") {
		t.Error("Output should contain another-opt test")
	}
}

func TestRoundTrip(t *testing.T) {
	input := `auto eth0
iface eth0 inet static
	address 192.168.1.10
	netmask 255.255.255.0
	gateway 192.168.1.1
	dns-nameservers 8.8.8.8
`

	file, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	var buf bytes.Buffer
	if err := GenerateTo(&buf, file); err != nil {
		t.Fatalf("GenerateTo failed: %v", err)
	}

	output := buf.String()

	if !strings.Contains(output, "auto eth0") {
		t.Error("Round trip failed: missing 'auto eth0'")
	}

	if !strings.Contains(output, "iface eth0 inet static") {
		t.Error("Round trip failed: missing 'iface eth0 inet static'")
	}

	if !strings.Contains(output, "address 192.168.1.10") {
		t.Error("Round trip failed: missing address")
	}
}

func intPtr(i int) *int {
	return &i
}

