package netinterface

import (
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	input := `# Example configuration
source /etc/network/interfaces.d/*

auto eth0
allow-hotplug eth0
iface eth0 inet static
	address 192.168.1.10
	netmask 255.255.255.0
	gateway 192.168.1.1
	dns-nameservers 8.8.8.8 8.8.4.4
	dns-search example.com

iface eth1 inet dhcp

iface br0 inet static
	bridge_ports eth0 eth1
	bridge_stp on

auto lo
iface lo inet loopback
`

	file, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(file.Sources) != 1 || file.Sources[0] != "/etc/network/interfaces.d/*" {
		t.Errorf("Expected 1 source, got %v", file.Sources)
	}

	if len(file.Interfaces) != 4 {
		t.Errorf("Expected 4 interfaces, got %d", len(file.Interfaces))
	}

	iface0 := file.Interfaces[0]
	if iface0.Name != "eth0" {
		t.Errorf("Expected name eth0, got %s", iface0.Name)
	}
	if !iface0.Auto {
		t.Error("Expected Auto=true for eth0")
	}
	if !iface0.Hotplug {
		t.Error("Expected Hotplug=true for eth0")
	}
	if iface0.Family != "inet" {
		t.Errorf("Expected family inet, got %s", iface0.Family)
	}
	if iface0.Method != "static" {
		t.Errorf("Expected method static, got %s", iface0.Method)
	}
	if iface0.Address != "192.168.1.10" {
		t.Errorf("Expected address 192.168.1.10, got %s", iface0.Address)
	}
	if iface0.Netmask != "255.255.255.0" {
		t.Errorf("Expected netmask 255.255.255.0, got %s", iface0.Netmask)
	}
	if iface0.Gateway != "192.168.1.1" {
		t.Errorf("Expected gateway 192.168.1.1, got %s", iface0.Gateway)
	}
	if len(iface0.DNS) != 2 {
		t.Errorf("Expected 2 DNS servers, got %d", len(iface0.DNS))
	}
	if len(iface0.DNSSearch) != 1 || iface0.DNSSearch[0] != "example.com" {
		t.Errorf("Expected DNS search example.com, got %v", iface0.DNSSearch)
	}

	iface1 := file.Interfaces[1]
	if iface1.Name != "eth1" || iface1.Method != "dhcp" {
		t.Errorf("Expected eth1 dhcp, got %s %s", iface1.Name, iface1.Method)
	}

	iface3 := file.Interfaces[3]
	if iface3.Name != "lo" || iface3.Method != "loopback" || !iface3.Auto {
		t.Errorf("Expected lo loopback auto, got %s %s auto=%v", iface3.Name, iface3.Method, iface3.Auto)
	}
}

func TestParseHwAddress(t *testing.T) {
	input := `iface eth0 inet static
	hwaddress ethernet 00:11:22:33:44:55
	address 192.168.1.10
	netmask 255.255.255.0
`

	file, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if file.Interfaces[0].HwAddress != "ethernet 00:11:22:33:44:55" {
		t.Errorf("Expected hwaddress 'ethernet 00:11:22:33:44:55', got '%s'", file.Interfaces[0].HwAddress)
	}
}

func TestParseCustom(t *testing.T) {
	input := `iface eth0 inet static
	address 192.168.1.10
	netmask 255.255.255.0
	custom-option value1
	custom-option value2
	another-custom test
`

	file, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	custom := file.Interfaces[0].Custom
	if len(custom["custom-option"]) != 2 {
		t.Errorf("Expected 2 values for custom-option, got %d", len(custom["custom-option"]))
	}
	if custom["another-custom"][0] != "test" {
		t.Errorf("Expected 'test' for another-custom, got '%s'", custom["another-custom"][0])
	}
}

func TestParseScripts(t *testing.T) {
	input := `iface eth0 inet static
	address 192.168.1.10
	pre-up /bin/do-something
	up /bin/do-other
	post-up echo "done"
	pre-down /bin/cleanup
	down /bin/stop
	post-down echo "cleaned"
`

	file, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	iface := file.Interfaces[0]
	if len(iface.PreUp) != 1 || iface.PreUp[0] != "/bin/do-something" {
		t.Errorf("Expected PreUp '/bin/do-something', got %v", iface.PreUp)
	}
	if len(iface.Up) != 1 || iface.Up[0] != "/bin/do-other" {
		t.Errorf("Expected Up '/bin/do-other', got %v", iface.Up)
	}
	if len(iface.PostUp) != 1 || iface.PostUp[0] != "echo \"done\"" {
		t.Errorf("Expected PostUp 'echo \"done\"', got %v", iface.PostUp)
	}
}

func TestParseEmpty(t *testing.T) {
	file, err := Parse(strings.NewReader(""))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if file == nil {
		t.Fatal("File should not be nil")
	}
	if len(file.Sources) != 0 || len(file.Interfaces) != 0 {
		t.Errorf("Expected empty file, got sources=%d interfaces=%d", len(file.Sources), len(file.Interfaces))
	}
}

