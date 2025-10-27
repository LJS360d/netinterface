package netinterface

import (
	"net"
)

type File struct {
	Sources    []string
	Interfaces []Interface
}

type Interface struct {
	Name       string
	Auto       bool
	Hotplug    bool
	Family     string
	Method     string
	Address    string
	Netmask    string
	Network    string
	Broadcast  string
	Gateway    string
	DNS        []net.IP
	DNSSearch  []string
	HwAddress  string
	PreUp      []string
	Up         []string
	PostUp     []string
	PreDown    []string
	Down       []string
	PostDown   []string
	Metric     *int
	MTU        *int
	BridgeOpts BridgeOptions
	Wireless   WirelessOptions
	Addresses  []AddressConfig
	Custom     map[string][]string
}

type AddressConfig struct {
	Address net.IP
	Netmask net.IP
}

type BridgeOptions struct {
	Ports     []string
	Stp       bool
	Ageing    int
	HelloTime int
	MaxAge    int
	Fd        bool
}

type WirelessOptions struct {
	SSID     string
	Mode     string
	Security string
	Key      string
}

