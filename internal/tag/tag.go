// Package tag provides port tagging based on well-known service mappings.
package tag

import "fmt"

// Tag represents a human-readable label for a port/protocol pair.
type Tag struct {
	Port     int
	Protocol string
	Label    string
}

// Tagger maps ports to known service tags.
type Tagger struct {
	custom map[string]string
}

// wellKnown contains common port-to-service mappings.
var wellKnown = map[string]string{
	"tcp:22":   "ssh",
	"tcp:80":   "http",
	"tcp:443":  "https",
	"tcp:3306": "mysql",
	"tcp:5432": "postgres",
	"tcp:6379": "redis",
	"tcp:8080": "http-alt",
	"udp:53":   "dns",
	"udp:123":  "ntp",
}

// New returns a Tagger with optional custom mappings merged over well-known ones.
func New(custom map[string]string) *Tagger {
	merged := make(map[string]string, len(wellKnown)+len(custom))
	for k, v := range wellKnown {
		merged[k] = v
	}
	for k, v := range custom {
		merged[k] = v
	}
	return &Tagger{custom: merged}
}

// Lookup returns the Tag for a given port and protocol.
// If no mapping exists, the label is "unknown".
func (t *Tagger) Lookup(port int, protocol string) Tag {
	key := fmt.Sprintf("%s:%d", protocol, port)
	label, ok := t.custom[key]
	if !ok {
		label = "unknown"
	}
	return Tag{Port: port, Protocol: protocol, Label: label}
}

// LookupAll returns tags for a slice of port/protocol pairs.
func (t *Tagger) LookupAll(ports []int, protocol string) []Tag {
	tags := make([]Tag, 0, len(ports))
	for _, p := range ports {
		tags = append(tags, t.Lookup(p, protocol))
	}
	return tags
}
