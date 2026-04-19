// Package fingerprint generates stable identifiers for port events.
package fingerprint

import (
	"crypto/sha256"
	"fmt"
	"strings"

	"github.com/user/portwatch/internal/alert"
)

// Fingerprint is a stable hash identifying a unique port event signature.
type Fingerprint string

// Generator produces fingerprints from alert events.
type Generator struct {
	fields []string
}

// New returns a Generator using the given fields for hashing.
// Valid fields: "port", "proto", "process", "action".
func New(fields []string) *Generator {
	if len(fields) == 0 {
		fields = []string{"port", "proto", "action"}
	}
	return &Generator{fields: fields}
}

// Compute returns a Fingerprint for the given event.
func (g *Generator) Compute(e alert.Event) Fingerprint {
	parts := make([]string, 0, len(g.fields))
	for _, f := range g.fields {
		switch f {
		case "port":
			parts = append(parts, fmt.Sprintf("%d", e.Port.Port))
		case "proto":
			parts = append(parts, e.Port.Proto)
		case "process":
			parts = append(parts, e.Port.Process)
		case "action":
			parts = append(parts, e.Action)
		}
	}
	raw := strings.Join(parts, ":")
	sum := sha256.Sum256([]byte(raw))
	return Fingerprint(fmt.Sprintf("%x", sum[:8]))
}

// ComputeAll returns fingerprints for a slice of events.
func (g *Generator) ComputeAll(events []alert.Event) []Fingerprint {
	out := make([]Fingerprint, len(events))
	for i, e := range events {
		out[i] = g.Compute(e)
	}
	return out
}
