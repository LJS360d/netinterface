package netinterface

import (
	"testing"
)

func TestValidateOK(t *testing.T) {
	file := &File{
		Interfaces: []Interface{
			{
				Name:   "eth0",
				Family: "inet",
				Method: "static",
				Address: "192.168.1.10",
				Netmask: "255.255.255.0",
				Gateway: "192.168.1.1",
			},
		},
	}

	errors := Validate(file)
	if len(errors) > 0 {
		t.Errorf("Expected no validation errors, got %v", errors)
	}
}

func TestValidateNameEmpty(t *testing.T) {
	file := &File{
		Interfaces: []Interface{
			{
				Family: "inet",
				Method: "static",
			},
		},
	}

	errors := Validate(file)
	if len(errors) == 0 {
		t.Error("Expected validation error for empty name")
	}

	found := false
	for _, err := range errors {
		if err.Field == "name" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected error with field='name'")
	}
}

func TestValidateInvalidFamily(t *testing.T) {
	file := &File{
		Interfaces: []Interface{
			{
				Name:   "eth0",
				Family: "invalid",
				Method: "static",
			},
		},
	}

	errors := Validate(file)
	if len(errors) == 0 {
		t.Error("Expected validation error for invalid family")
	}
}

func TestValidateInvalidAddress(t *testing.T) {
	file := &File{
		Interfaces: []Interface{
			{
				Name:    "eth0",
				Family:  "inet",
				Method:  "static",
				Address: "not-an-ip",
			},
		},
	}

	errors := Validate(file)
	if len(errors) == 0 {
		t.Error("Expected validation error for invalid address")
	}

	found := false
	for _, err := range errors {
		if err.Field == "address" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected error with field='address'")
	}
}

func TestValidateInvalidMTU(t *testing.T) {
	tooLow := 50
	tooHigh := 65536
	file := &File{
		Interfaces: []Interface{
			{
				Name:   "eth0",
				Family: "inet",
				Method: "static",
				MTU:    &tooLow,
			},
			{
				Name:   "eth1",
				Family: "inet",
				Method: "static",
				MTU:    &tooHigh,
			},
		},
	}

	errors := Validate(file)
	if len(errors) < 2 {
		t.Errorf("Expected at least 2 validation errors, got %d", len(errors))
	}

	found := false
	for _, err := range errors {
		if err.Field == "mtu" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected error with field='mtu'")
	}
}

func TestValidateDuplicateInterface(t *testing.T) {
	file := &File{
		Interfaces: []Interface{
			{
				Name:   "eth0",
				Family: "inet",
				Method: "static",
			},
			{
				Name:   "eth0",
				Family: "inet",
				Method: "dhcp",
			},
		},
	}

	errors := Validate(file)
	if len(errors) == 0 {
		t.Error("Expected validation error for duplicate interface names")
	}
}

func TestValidateAndReport(t *testing.T) {
	file := &File{
		Interfaces: []Interface{
			{
				Family: "inet",
				Method: "static",
			},
		},
	}

	err := ValidateAndReport(file)
	if err == nil {
		t.Error("Expected validation error")
	}
}

func TestValidateOKAndReport(t *testing.T) {
	file := &File{
		Interfaces: []Interface{
			{
				Name:   "eth0",
				Family: "inet",
				Method: "static",
			},
		},
	}

	err := ValidateAndReport(file)
	if err != nil {
		t.Errorf("Expected no validation error, got %v", err)
	}
}

