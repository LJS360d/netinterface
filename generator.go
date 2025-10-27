package netinterface

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func Generate(f *File, path string) error {
	w, err := os.Create(path)
	if err != nil {
		return err
	}
	defer w.Close()
	return GenerateTo(w, f)
}

func GenerateTo(w io.Writer, f *File) error {
	buf := bufio.NewWriter(w)
	defer buf.Flush()

	for _, source := range f.Sources {
		if _, err := fmt.Fprintf(buf, "source %s\n", source); err != nil {
			return err
		}
	}

	if len(f.Sources) > 0 {
		if _, err := fmt.Fprintln(buf); err != nil {
			return err
		}
	}

	for _, iface := range f.Interfaces {
		if err := writeInterface(buf, iface); err != nil {
			return err
		}
	}

	return nil
}

func writeInterface(buf *bufio.Writer, iface Interface) error {
	if iface.Auto {
		if _, err := fmt.Fprintf(buf, "auto %s\n", iface.Name); err != nil {
			return err
		}
	}

	if iface.Hotplug {
		if _, err := fmt.Fprintf(buf, "allow-hotplug %s\n", iface.Name); err != nil {
			return err
		}
	}

	if iface.Name != "" && iface.Family != "" && iface.Method != "" {
		if _, err := fmt.Fprintf(buf, "iface %s %s %s\n", iface.Name, iface.Family, iface.Method); err != nil {
			return err
		}

		if iface.Address != "" {
			if _, err := fmt.Fprintf(buf, "\taddress %s\n", iface.Address); err != nil {
				return err
			}
		}

		if iface.Netmask != "" {
			if _, err := fmt.Fprintf(buf, "\tnetmask %s\n", iface.Netmask); err != nil {
				return err
			}
		}

		if iface.Network != "" {
			if _, err := fmt.Fprintf(buf, "\tnetwork %s\n", iface.Network); err != nil {
				return err
			}
		}

		if iface.Broadcast != "" {
			if _, err := fmt.Fprintf(buf, "\tbroadcast %s\n", iface.Broadcast); err != nil {
				return err
			}
		}

		if iface.Gateway != "" {
			if _, err := fmt.Fprintf(buf, "\tgateway %s\n", iface.Gateway); err != nil {
				return err
			}
		}

		if iface.HwAddress != "" {
			if _, err := fmt.Fprintf(buf, "\thwaddress %s\n", iface.HwAddress); err != nil {
				return err
			}
		}

		if iface.Metric != nil {
			if _, err := fmt.Fprintf(buf, "\tmetric %d\n", *iface.Metric); err != nil {
				return err
			}
		}

		if iface.MTU != nil {
			if _, err := fmt.Fprintf(buf, "\tmtu %d\n", *iface.MTU); err != nil {
				return err
			}
		}

		if len(iface.DNS) > 0 {
			ips := make([]string, 0, len(iface.DNS))
			for _, ip := range iface.DNS {
				ips = append(ips, ip.String())
			}
			if _, err := fmt.Fprintf(buf, "\tdns-nameservers %s\n", strings.Join(ips, " ")); err != nil {
				return err
			}
		}

		if len(iface.DNSSearch) > 0 {
			if _, err := fmt.Fprintf(buf, "\tdns-search %s\n", strings.Join(iface.DNSSearch, " ")); err != nil {
				return err
			}
		}

		for _, cmd := range iface.PreUp {
			if _, err := fmt.Fprintf(buf, "\tpre-up %s\n", cmd); err != nil {
				return err
			}
		}

		for _, cmd := range iface.Up {
			if _, err := fmt.Fprintf(buf, "\tup %s\n", cmd); err != nil {
				return err
			}
		}

		for _, cmd := range iface.PostUp {
			if _, err := fmt.Fprintf(buf, "\tpost-up %s\n", cmd); err != nil {
				return err
			}
		}

		for _, cmd := range iface.PreDown {
			if _, err := fmt.Fprintf(buf, "\tpre-down %s\n", cmd); err != nil {
				return err
			}
		}

		for _, cmd := range iface.Down {
			if _, err := fmt.Fprintf(buf, "\tdown %s\n", cmd); err != nil {
				return err
			}
		}

		for _, cmd := range iface.PostDown {
			if _, err := fmt.Fprintf(buf, "\tpost-down %s\n", cmd); err != nil {
				return err
			}
		}

		if len(iface.BridgeOpts.Ports) > 0 {
			if _, err := fmt.Fprintf(buf, "\tbridge_ports %s\n", strings.Join(iface.BridgeOpts.Ports, " ")); err != nil {
				return err
			}
		}

		if iface.BridgeOpts.Stp {
			if _, err := fmt.Fprintf(buf, "\tbridge_stp on\n"); err != nil {
				return err
			}
		}

		if iface.BridgeOpts.Ageing > 0 {
			if _, err := fmt.Fprintf(buf, "\tbridge_ageing %d\n", iface.BridgeOpts.Ageing); err != nil {
				return err
			}
		}

		if iface.Wireless.SSID != "" {
			if _, err := fmt.Fprintf(buf, "\twireless-essid %s\n", iface.Wireless.SSID); err != nil {
				return err
			}
		}

		if iface.Wireless.Mode != "" {
			if _, err := fmt.Fprintf(buf, "\twireless-mode %s\n", iface.Wireless.Mode); err != nil {
				return err
			}
		}

		if iface.Wireless.Key != "" {
			if _, err := fmt.Fprintf(buf, "\twireless-key %s\n", iface.Wireless.Key); err != nil {
				return err
			}
		}

		if iface.Custom != nil {
			for key, values := range iface.Custom {
				for _, value := range values {
					if _, err := fmt.Fprintf(buf, "\t%s %s\n", key, value); err != nil {
						return err
					}
				}
			}
		}
	}

	if _, err := fmt.Fprintln(buf); err != nil {
		return err
	}

	return nil
}

