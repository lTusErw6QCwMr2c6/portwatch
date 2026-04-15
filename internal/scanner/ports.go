package scanner

import (
	"bufio"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// PortKey uniquely identifies a port + protocol combination.
type PortKey struct {
	Port  int    `json:"port"`
	Proto string `json:"proto"`
}

// DiffResult holds ports that were opened or closed between two scans.
type DiffResult struct {
	Opened []PortKey
	Closed []PortKey
}

// ScanPorts returns the set of open ports in [start, end] using ss/netstat.
func ScanPorts(start, end int) (map[PortKey]bool, error) {
	if start < 0 || end > 65535 || start > end {
		return nil, fmt.Errorf("invalid port range: %d-%d", start, end)
	}
	out, err := exec.Command("ss", "-tlunp").Output()
	if err != nil {
		return nil, fmt.Errorf("ss command failed: %w", err)
	}
	return parseSS(out, start, end), nil
}

func parseSS(data []byte, start, end int) map[PortKey]bool {
	result := make(map[PortKey]bool)
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 5 {
			continue
		}
		proto := strings.ToLower(fields[0])
		if proto != "tcp" && proto != "udp" {
			continue
		}
		addr := fields[4]
		port := extractPort(addr)
		if port < start || port > end {
			continue
		}
		result[PortKey{Port: port, Proto: proto}] = true
	}
	return result
}

func extractPort(addr string) int {
	if i := strings.LastIndex(addr, ":"); i >= 0 {
		p, err := strconv.Atoi(addr[i+1:])
		if err == nil {
			return p
		}
	}
	return -1
}

// Diff computes the difference between two port sets.
func Diff(prev, curr map[PortKey]bool) DiffResult {
	var result DiffResult
	for k := range curr {
		if !prev[k] {
			result.Opened = append(result.Opened, k)
		}
	}
	for k := range prev {
		if !curr[k] {
			result.Closed = append(result.Closed, k)
		}
	}
	return result
}

// FormatPort returns a human-readable port string.
func FormatPort(k PortKey) string {
	return fmt.Sprintf("%s/%d", k.Proto, k.Port)
}
