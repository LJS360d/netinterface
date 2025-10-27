package netinterface

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
)

type state int

const (
	stateGlobal state = iota
	stateInIface
)

type token struct {
	line   int
	fields []string
}

func tokenizeLine(line string, lineno int) token {
	fields := strings.Fields(line)
	if len(fields) == 0 {
		fields = []string{""}
	}
	return token{line: lineno, fields: fields}
}

func Parse(r io.Reader) (*File, error) {
	file := &File{
		Sources:    []string{},
		Interfaces: []Interface{},
	}

	scanner := bufio.NewScanner(r)
	var lineno int
	var current *Interface
	var state state = stateGlobal

	for scanner.Scan() {
		lineno++
		line := strings.TrimSpace(scanner.Text())

		tok := tokenizeLine(line, lineno)

		if len(tok.fields) == 0 || tok.fields[0] == "" {
			if current != nil && state == stateInIface {
				state = stateGlobal
			}
			continue
		}

		if tok.fields[0][0] == '#' {
			continue
		}

		switch tok.fields[0] {
		case "source", "source-directory":
			if len(tok.fields) > 1 {
				file.Sources = append(file.Sources, tok.fields[1])
			}

		case "auto":
			if len(tok.fields) > 1 {
				for _, name := range tok.fields[1:] {
					iface := findOrCreateInterface(file, name)
					iface.Auto = true
				}
			}

		case "allow-hotplug":
			if len(tok.fields) > 1 {
				for _, name := range tok.fields[1:] {
					iface := findOrCreateInterface(file, name)
					iface.Hotplug = true
				}
			}

		case "iface":
			if len(tok.fields) >= 3 {
				state = stateInIface
				name := tok.fields[1]
				family := tok.fields[2]
				method := ""
				if len(tok.fields) > 3 {
					method = tok.fields[3]
				}

				current = findOrCreateInterface(file, name)
				current.Family = family
				current.Method = method
			}

		case "address":
			if current != nil && len(tok.fields) > 1 {
				current.Address = tok.fields[1]
			}

		case "netmask":
			if current != nil && len(tok.fields) > 1 {
				current.Netmask = tok.fields[1]
			}

		case "network":
			if current != nil && len(tok.fields) > 1 {
				current.Network = tok.fields[1]
			}

		case "broadcast":
			if current != nil && len(tok.fields) > 1 {
				current.Broadcast = tok.fields[1]
			}

		case "gateway":
			if current != nil && len(tok.fields) > 1 {
				current.Gateway = tok.fields[1]
			}

		case "metric":
			if current != nil && len(tok.fields) > 1 {
				var metric int
				fmt.Sscanf(tok.fields[1], "%d", &metric)
				current.Metric = &metric
			}

		case "mtu":
			if current != nil && len(tok.fields) > 1 {
				var mtu int
				fmt.Sscanf(tok.fields[1], "%d", &mtu)
				current.MTU = &mtu
			}

		case "dns-nameservers":
			if current != nil && len(tok.fields) > 1 {
				for _, ns := range tok.fields[1:] {
					if ip := net.ParseIP(ns); ip != nil {
						current.DNS = append(current.DNS, ip)
					}
				}
			}

		case "dns-search":
			if current != nil && len(tok.fields) > 1 {
				current.DNSSearch = append(current.DNSSearch, tok.fields[1:]...)
			}

		case "hwaddress":
			if current != nil && len(tok.fields) > 1 {
				current.HwAddress = strings.Join(tok.fields[1:], " ")
			}

		case "pre-up":
			if current != nil && len(tok.fields) > 1 {
				current.PreUp = append(current.PreUp, strings.Join(tok.fields[1:], " "))
			}

		case "up":
			if current != nil && len(tok.fields) > 1 {
				current.Up = append(current.Up, strings.Join(tok.fields[1:], " "))
			}

		case "post-up":
			if current != nil && len(tok.fields) > 1 {
				current.PostUp = append(current.PostUp, strings.Join(tok.fields[1:], " "))
			}

		case "pre-down":
			if current != nil && len(tok.fields) > 1 {
				current.PreDown = append(current.PreDown, strings.Join(tok.fields[1:], " "))
			}

		case "down":
			if current != nil && len(tok.fields) > 1 {
				current.Down = append(current.Down, strings.Join(tok.fields[1:], " "))
			}

		case "post-down":
			if current != nil && len(tok.fields) > 1 {
				current.PostDown = append(current.PostDown, strings.Join(tok.fields[1:], " "))
			}

		case "bridge_ports":
			if current != nil && len(tok.fields) > 1 {
				current.BridgeOpts.Ports = tok.fields[1:]
			}

		case "bridge_stp":
			if current != nil && len(tok.fields) > 1 {
				current.BridgeOpts.Stp = tok.fields[1] == "on" || tok.fields[1] == "yes"
			}

		case "bridge_ageing":
			if current != nil && len(tok.fields) > 1 {
				fmt.Sscanf(tok.fields[1], "%d", &current.BridgeOpts.Ageing)
			}

		case "wireless-essid":
			if current != nil && len(tok.fields) > 1 {
				current.Wireless.SSID = tok.fields[1]
			}

		case "wireless-mode":
			if current != nil && len(tok.fields) > 1 {
				current.Wireless.Mode = tok.fields[1]
			}

		case "wireless-key":
			if current != nil && len(tok.fields) > 1 {
				current.Wireless.Key = strings.Join(tok.fields[1:], " ")
			}

		default:
			if current != nil && tok.fields[0] != "" {
				parseCustomOption(current, tok.fields[0], strings.Join(tok.fields[1:], " "))
			}
		}
	}

	return file, scanner.Err()
}

func findOrCreateInterface(file *File, name string) *Interface {
	for i := range file.Interfaces {
		if file.Interfaces[i].Name == name {
			if file.Interfaces[i].Custom == nil {
				file.Interfaces[i].Custom = make(map[string][]string)
			}
			return &file.Interfaces[i]
		}
	}
	iface := Interface{
		Name:   name,
		Custom: make(map[string][]string),
	}
	file.Interfaces = append(file.Interfaces, iface)
	return &file.Interfaces[len(file.Interfaces)-1]
}

func parseCustomOption(iface *Interface, key, value string) {
	value = strings.TrimSpace(value)
	switch {
	case strings.HasPrefix(key, "pre-up"):
		iface.PreUp = append(iface.PreUp, value)
	case strings.HasPrefix(key, "up"):
		iface.Up = append(iface.Up, value)
	case strings.HasPrefix(key, "post-up"):
		iface.PostUp = append(iface.PostUp, value)
	case strings.HasPrefix(key, "pre-down"):
		iface.PreDown = append(iface.PreDown, value)
	case strings.HasPrefix(key, "down"):
		iface.Down = append(iface.Down, value)
	case strings.HasPrefix(key, "post-down"):
		iface.PostDown = append(iface.PostDown, value)
	default:
		if iface.Custom == nil {
			iface.Custom = make(map[string][]string)
		}
		if value != "" {
			iface.Custom[key] = append(iface.Custom[key], value)
		}
	}
}

