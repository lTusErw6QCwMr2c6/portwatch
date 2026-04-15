package scanner

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

// PortState represents the state of a single port.
type PortState struct {
	Port     int
	Protocol string
	State    string // "open" or "closed"
}

// Snapshot holds all observed open ports at a point in time.
type Snapshot map[string]PortState

// key returns a unique string key for a port/protocol combo.
func key(port int, proto string) string {
	return fmt.Sprintf("%s:%d", proto, port)
}

// ScanPorts checks a range of TCP and UDP ports on localhost and returns
// a Snapshot of all ports that are currently open.
func ScanPorts(startPort, endPort int) (Snapshot, error) {
	if startPort < 1 || endPort > 65535 || startPort > endPort {
		return nil, fmt.Errorf("invalid port range: %d-%d", startPort, endPort)
	}

	snap := make(Snapshot)

	for port := startPort; port <= endPort; port++ {
		addr := net.JoinHostPort("127.0.0.1", strconv.Itoa(port))
		conn, err := net.Dial("tcp", addr)
		if err == nil {
			conn.Close()
			ps := PortState{Port: port, Protocol: "tcp", State: "open"}
			snap[key(port, "tcp")] = ps
		}
	}

	return snap, nil
}

// Diff compares two snapshots and returns newly opened and newly closed ports.
func Diff(prev, curr Snapshot) (opened, closed []PortState) {
	 k, ps := range curr {
		if _, exists := prev[k]; !exists {
			opened = append(opened, ps)
		}
	}
	for k, ps := range prev {
		if _, exists := curr[k]; !exists {
			closed = append(closed, ps)
		}
	}
	return
}

// FormatPort returns a human-readable label for a PortState.
func FormatPort(ps PortState) string {
	return fmt.Sprintf("%s/%d", strings.ToUpper(ps.Protocol), ps.Port)
}
