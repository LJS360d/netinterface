package main

import (
	"fmt"
	"net"
	"os"

	"github.com/LJS360d/netinterface"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <command>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Commands:\n")
		fmt.Fprintf(os.Stderr, "  read    - Read and display /etc/network/interfaces\n")
		fmt.Fprintf(os.Stderr, "  write   - Create a sample interfaces file\n")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "read":
		readExample()
	case "write":
		writeExample()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}

func readExample() {
	file, err := netinterface.ReadDefault()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("=== Sources ===")
	for _, source := range file.Sources {
		fmt.Printf("  source %s\n", source)
	}

	fmt.Println("\n=== Interfaces ===")
	for _, iface := range file.Interfaces {
		fmt.Printf("Interface: %s\n", iface.Name)
		if iface.Auto {
			fmt.Println("  Auto: enabled")
		}
		if iface.Hotplug {
			fmt.Println("  Hotplug: enabled")
		}
		fmt.Printf("  Family: %s\n", iface.Family)
		fmt.Printf("  Method: %s\n", iface.Method)
		if iface.Address != "" {
			fmt.Printf("  Address: %s\n", iface.Address)
		}
		if iface.Netmask != "" {
			fmt.Printf("  Netmask: %s\n", iface.Netmask)
		}
		if iface.Gateway != "" {
			fmt.Printf("  Gateway: %s\n", iface.Gateway)
		}
		if len(iface.DNS) > 0 {
			fmt.Printf("  DNS: ")
			for _, ip := range iface.DNS {
				fmt.Printf("%s ", ip)
			}
			fmt.Println()
		}
		fmt.Println()
	}
}

func writeExample() {
	file := &netinterface.File{
		Sources: []string{"/etc/network/interfaces.d/*"},
		Interfaces: []netinterface.Interface{
			{
				Name:   "lo",
				Auto:   true,
				Family: "inet",
				Method: "loopback",
			},
			{
				Name:     "eth0",
				Auto:     true,
				Hotplug:  true,
				Family:   "inet",
				Method:   "static",
				Address:  "192.168.1.100",
				Netmask:  "255.255.255.0",
				Gateway:  "192.168.1.1",
				DNS:      []net.IP{net.ParseIP("8.8.8.8"), net.ParseIP("8.8.4.4")},
				DNSSearch: []string{"local"},
			},
		},
	}

	if err := netinterface.WriteFile("sample_interfaces", file); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing file: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Sample interfaces file written to 'sample_interfaces'")
}

