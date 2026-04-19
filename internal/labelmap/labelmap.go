// Package labelmap attaches arbitrary string labels to events based on port or protocol rules.
package labelmap

import (
	"sync"

	"github.com/user/portwatch/internal/alert"
)

// Rule maps a port and optional protocol to a set of labels.
type Rule struct {
	Port     int
	Protocol string // empty matches any
	Labels   []string
}

// Mapper holds a set of labelling rules.
type Mapper struct {
	mu    sync.RWMutex
	rules []Rule
}

// New returns a Mapper pre-loaded with the given rules.
func New(rules []Rule) *Mapper {
	return &Mapper{rules: rules}
}

// Add appends a rule to the mapper.
func (m *Mapper) Add(r Rule) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.rules = append(m.rules, r)
}

// Apply returns all labels that match the event's port and protocol.
func (m *Mapper) Apply(e alert.Event) []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	seen := map[string]struct{}{}
	var out []string
	for _, r := range m.rules {
		if r.Port != e.Port.Number {
			continue
		}
		if r.Protocol != "" && r.Protocol != e.Port.Protocol {
			continue
		}
		for _, l := range r.Labels {
			if _, ok := seen[l]; !ok {
				seen[l] = struct{}{}
				out = append(out, l)
			}
		}
	}
	return out
}

// Reset removes all rules.
func (m *Mapper) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.rules = nil
}
