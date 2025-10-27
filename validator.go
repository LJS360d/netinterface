package netinterface

import (
	"fmt"
	"net"
	"strings"
)

type ValidationError struct {
	Interface string
	Field     string
	Message   string
}

func (e ValidationError) Error() string {
	if e.Interface != "" {
		return fmt.Sprintf("interface %s: %s: %s", e.Interface, e.Field, e.Message)
	}
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

func Validate(f *File) []ValidationError {
	var errors []ValidationError

	for _, source := range f.Sources {
		if source == "" {
			errors = append(errors, ValidationError{Field: "source", Message: "empty source path"})
		}
	}

	interfaceNames := make(map[string]bool)
	for i := range f.Interfaces {
		iface := &f.Interfaces[i]
		if errs := validateInterface(iface); len(errs) > 0 {
			errors = append(errors, errs...)
		}

		if interfaceNames[iface.Name] {
			errors = append(errors, ValidationError{
				Interface: iface.Name,
				Field:     "name",
				Message:   "duplicate interface name",
			})
		}
		interfaceNames[iface.Name] = true
	}

	return errors
}

func validateInterface(iface *Interface) []ValidationError {
	var errors []ValidationError

	if iface.Name == "" {
		errors = append(errors, ValidationError{
			Interface: iface.Name,
			Field:     "name",
			Message:   "interface name cannot be empty",
		})
	}

	if iface.Family != "" && iface.Family != "inet" && iface.Family != "inet6" {
		errors = append(errors, ValidationError{
			Interface: iface.Name,
			Field:     "family",
			Message:   fmt.Sprintf("invalid family: %s (must be 'inet' or 'inet6')", iface.Family),
		})
	}

	if iface.Method == "" && (iface.Name != "" || iface.Family != "") {
		errors = append(errors, ValidationError{
			Interface: iface.Name,
			Field:     "method",
			Message:   "method cannot be empty",
		})
	}

	if iface.Address != "" {
		ip := net.ParseIP(iface.Address)
		if ip == nil {
			errors = append(errors, ValidationError{
				Interface: iface.Name,
				Field:     "address",
				Message:   fmt.Sprintf("invalid IP address: %s", iface.Address),
			})
		}
	}

	if iface.Netmask != "" {
		ip := net.ParseIP(iface.Netmask)
		if ip == nil {
			errors = append(errors, ValidationError{
				Interface: iface.Name,
				Field:     "netmask",
				Message:   fmt.Sprintf("invalid IP address: %s", iface.Netmask),
			})
		}
	}

	if iface.Gateway != "" {
		ip := net.ParseIP(iface.Gateway)
		if ip == nil {
			errors = append(errors, ValidationError{
				Interface: iface.Name,
				Field:     "gateway",
				Message:   fmt.Sprintf("invalid IP address: %s", iface.Gateway),
			})
		}
	}

	if iface.Network != "" {
		ip := net.ParseIP(iface.Network)
		if ip == nil {
			errors = append(errors, ValidationError{
				Interface: iface.Name,
				Field:     "network",
				Message:   fmt.Sprintf("invalid IP address: %s", iface.Network),
			})
		}
	}

	if iface.Broadcast != "" {
		ip := net.ParseIP(iface.Broadcast)
		if ip == nil {
			errors = append(errors, ValidationError{
				Interface: iface.Name,
				Field:     "broadcast",
				Message:   fmt.Sprintf("invalid IP address: %s", iface.Broadcast),
			})
		}
	}

	if iface.Metric != nil && *iface.Metric < 0 {
		errors = append(errors, ValidationError{
			Interface: iface.Name,
			Field:     "metric",
			Message:   "metric cannot be negative",
		})
	}

	if iface.MTU != nil && (*iface.MTU < 68 || *iface.MTU > 65535) {
		errors = append(errors, ValidationError{
			Interface: iface.Name,
			Field:     "mtu",
			Message:   fmt.Sprintf("MTU out of valid range (68-65535): %d", *iface.MTU),
		})
	}

	for _, ip := range iface.DNS {
		if ip == nil {
			errors = append(errors, ValidationError{
				Interface: iface.Name,
				Field:     "dns",
				Message:   "contains nil IP address",
			})
		}
	}

	if len(iface.BridgeOpts.Ports) > 0 && iface.Method != "static" && iface.Method != "dhcp" && iface.Method != "manual" {
		errors = append(errors, ValidationError{
			Interface: iface.Name,
			Field:     "bridge_ports",
			Message:   "bridge options only valid for bridge interfaces",
		})
	}

	if iface.Wireless.SSID != "" && iface.Method != "dhcp" && iface.Method != "static" && iface.Method != "manual" {
		errors = append(errors, ValidationError{
			Interface: iface.Name,
			Field:     "wireless",
			Message:   "wireless options only valid for wireless interfaces",
		})
	}

	return errors
}

func ValidateAndReport(f *File) error {
	errors := Validate(f)
	if len(errors) == 0 {
		return nil
	}

	messages := make([]string, len(errors))
	for i, err := range errors {
		messages[i] = err.Error()
	}

	return fmt.Errorf("validation failed:\n%s", strings.Join(messages, "\n"))
}

