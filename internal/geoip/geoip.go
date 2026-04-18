// Package geoip provides lightweight IP geolocation lookups for port events.
package geoip

import (
	"net"
	"strings"
)

// Info holds geolocation metadata for an IP address.
type Info struct {
	IP      string
	Country string
	Private bool
}

// Lookup resolves geolocation info for a given IP string.
// For private/loopback addresses it returns a Private=true result without
// an external call. Real deployments would integrate a MaxMind DB here.
func Lookup(ipStr string) Info {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return Info{IP: ipStr, Country: "unknown"}
	}
	if isPrivate(ip) {
		return Info{IP: ipStr, Country: "private", Private: true}
	}
	// Stub: real implementation would query a local GeoIP database.
	return Info{IP: ipStr, Country: "unknown"}
}

// String returns a human-readable summary of the geolocation info.
func (i Info) String() string {
	if i.Private {
		return i.IP + " (private)"
	}
	return i.IP + " (" + i.Country + ")"
}

var privateRanges = []string{
	"10.",
	"192.168.",
	"172.",
	"127.",
	"::1",
	"fc",
	"fd",
}

func isPrivate(ip net.IP) bool {
	s := ip.String()
	for _, prefix := range privateRanges {
		if strings.HasPrefix(s, prefix) {
			return true
		}
	}
	return false
}
