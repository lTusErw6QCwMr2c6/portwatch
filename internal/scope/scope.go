// Package scope restricts monitoring to a named set of ports or services.
package scope

import "sync"

// Entry describes a named scope entry.
type Entry struct {
	Name     string
	Ports    []uint16
	Protocol string // "tcp", "udp", or "" for both
}

// Scope holds a collection of named entries and provides membership checks.
type Scope struct {
	mu      sync.RWMutex
	entries map[string]Entry
}

// New returns an empty Scope.
func New() *Scope {
	return &Scope{entries: make(map[string]Entry)}
}

// Add registers a named entry, replacing any existing entry with the same name.
func (s *Scope) Add(e Entry) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries[e.Name] = e
}

// Remove deletes the entry with the given name.
func (s *Scope) Remove(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.entries, name)
}

// Match returns the name of the first entry whose port and protocol match,
// and true. Returns "", false if no entry matches.
func (s *Scope) Match(port uint16, proto string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, e := range s.entries {
		if e.Protocol != "" && e.Protocol != proto {
			continue
		}
		for _, p := range e.Ports {
			if p == port {
				return e.Name, true
			}
		}
	}
	return "", false
}

// Names returns all registered entry names.
func (s *Scope) Names() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]string, 0, len(s.entries))
	for k := range s.entries {
		out = append(out, k)
	}
	return out
}

// Len returns the number of registered entries.
func (s *Scope) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.entries)
}
