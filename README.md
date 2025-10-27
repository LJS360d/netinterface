# netinterface

A Go library for reading and writing Linux `/etc/network/interfaces` configuration files.

## Features

- **Complete support** for standard interfaces file directives
- **Type-safe** data structures for all configuration options
- **Validation** of IP addresses, MTU, and other fields
- **Parse and generate** interfaces files with proper formatting
- **Extensible** architecture for custom options

## Installation

```bash
go get github.com/LJS360d/netinterface
```

## Quick Start

### Most Common Use Case

For working with `/etc/network/interfaces` (99% of use cases):

```go
import "github.com/LJS360d/netinterface"

// Read
file, err := netinterface.ReadDefault()
if err != nil {
    panic(err)
}

// Modify
file.Interfaces[0].Address = "192.168.1.100"

// Write back
err := netinterface.WriteDefault(file)
```

That's it. The library uses `/etc/network/interfaces` by default.

### Reading Configuration

Read from a file:
```go
package main

import (
    "fmt"
    "github.com/LJS360d/netinterface"
)

func main() {
    file, err := netinterface.ReadFile("/etc/network/interfaces")
    if err != nil {
        panic(err)
    }

    for _, iface := range file.Interfaces {
        fmt.Printf("Interface: %s\n", iface.Name)
        fmt.Printf("  Family: %s\n", iface.Family)
        fmt.Printf("  Method: %s\n", iface.Method)
        fmt.Printf("  Address: %s\n", iface.Address)
        if iface.Auto {
            fmt.Println("  Auto: enabled")
        }
    }
}
```

Read from any `io.Reader`:
```go
import "io"

// From a string
reader := strings.NewReader("auto eth0\niface eth0 inet dhcp\n")
file, err := netinterface.Read(reader)
```

### Writing Configuration

Write to a file:
```go
package main

import (
    "net"
    "github.com/LJS360d/netinterface"
)

func main() {
    file := &netinterface.File{
        Interfaces: []netinterface.Interface{
            {
                Name:     "eth0",
                Auto:     true,
                Family:   "inet",
                Method:   "static",
                Address:  "192.168.1.10",
                Netmask:  "255.255.255.0",
                Gateway:  "192.168.1.1",
                DNS:      []net.IP{net.ParseIP("8.8.8.8"), net.ParseIP("8.8.4.4")},
            },
        },
    }

    if err := netinterface.WriteFile("/etc/network/interfaces", file); err != nil {
        panic(err)
    }
}
```

Write to any `io.Writer`:
```go
import "io"

var buf bytes.Buffer
err := netinterface.Write(&buf, file)
// Get the output as a string
output := buf.String()
```

## Supported Configuration Options

### Interface Configuration

- `address` - IP address
- `netmask` - Network mask
- `network` - Network address
- `broadcast` - Broadcast address
- `gateway` - Default gateway
- `hwaddress` - Hardware address
- `metric` - Route metric
- `mtu` - Maximum transmission unit

### DNS Configuration

- `dns-nameservers` - DNS servers
- `dns-search` - DNS search domains

### Bridge Configuration

- `bridge_ports` - Bridge port interfaces
- `bridge_stp` - Spanning tree protocol
- `bridge_ageing` - Ageing time

### Wireless Configuration

- `wireless-essid` - SSID
- `wireless-mode` - Operating mode
- `wireless-key` - Security key

### Script Hooks

- `pre-up`, `up`, `post-up` - Bring-up hooks
- `pre-down`, `down`, `post-down` - Bring-down hooks

### Custom Fields

For any options not explicitly supported, use the `Custom` map:

```go
iface.Custom = map[string][]string{
    "unknown-directive": []string{"value"},
}
```

Fields stored in `Custom` will be written back to the file on output.

## Validation

The library validates configuration before writing:

```go
errors := netinterface.Validate(file)
for _, err := range errors {
    fmt.Println(err)
}

// Or get a single combined error
if err := netinterface.ValidateAndReport(file); err != nil {
    log.Fatal(err)
}
```

## API Reference

### Types

#### `File`

Top-level container for the interfaces file.

```go
type File struct {
    Sources    []string
    Interfaces []Interface
}
```

#### `Interface`

Represents a single network interface configuration.

```go
type Interface struct {
    Name      string
    Auto      bool
    Hotplug   bool
    Family    string  // "inet" or "inet6"
    Method    string  // "dhcp", "static", "loopback", "manual", etc.
    Address   string
    Netmask   string
    Network   string
    Broadcast string
    Gateway   string
    DNS       []net.IP
    DNSSearch []string
    // ... and more
}
```

### Functions

#### `ReadDefault() (*File, error)`

Read from `/etc/network/interfaces` (most common case).

#### `Read(r io.Reader) (*File, error)`

Read an interfaces configuration from any `io.Reader`.

#### `ReadFile(path string) (*File, error)`

Read an interfaces file from disk.

#### `WriteDefault(f *File) error`

Write to `/etc/network/interfaces` (most common case).

#### `Write(w io.Writer, f *File) error`

Write an interfaces configuration to any `io.Writer` with validation.

#### `WriteFile(path string, f *File) error`

Write an interfaces file to disk with validation.

#### `Validate(file *File) []ValidationError`

Validate a configuration and return all errors.

#### `ValidateAndReport(file *File) error`

Validate a configuration and return a combined error.

#### Constants

- `DefaultInterfacesPath` - Default path `/etc/network/interfaces`

## Example Use Cases

### DHCP Interface

```go
iface := netinterface.Interface{
    Name:   "eth0",
    Auto:   true,
    Family: "inet",
    Method: "dhcp",
}
```

### Static IP with DNS and Hardware Address

```go
iface := netinterface.Interface{
    Name:      "eth0",
    Auto:      true,
    Family:    "inet",
    Method:    "static",
    Address:   "192.168.1.100",
    Netmask:   "255.255.255.0",
    Gateway:   "192.168.1.1",
    HwAddress: "ethernet 00:11:22:33:44:55",
    DNS:       []net.IP{net.ParseIP("8.8.8.8")},
    DNSSearch: []string{"local"},
}
```

### Bridge Interface

```go
iface := netinterface.Interface{
    Name:   "br0",
    Auto:   true,
    Family: "inet",
    Method: "static",
    BridgeOpts: netinterface.BridgeOptions{
        Ports: []string{"eth0", "eth1"},
        Stp:   true,
    },
}
```

### Interface with Hooks

```go
iface := netinterface.Interface{
    Name:   "eth0",
    Auto:   true,
    Family: "inet",
    Method: "static",
    PreUp:  []string{"ifconfig $IFACE promisc"},
    Up:     []string{"iptables -A INPUT -i $IFACE -j ACCEPT"},
    PostUp: []string{"logger Interface $IFACE is up"},
}
```

### Interface with Custom Fields

```go
iface := netinterface.Interface{
    Name:   "eth0",
    Family: "inet",
    Method: "static",
    Custom: map[string][]string{
        "pre-custom": []string{"do something"},
        "post-custom": []string{"do something else"},
    },
}
```

## License

See LICENSE file.
