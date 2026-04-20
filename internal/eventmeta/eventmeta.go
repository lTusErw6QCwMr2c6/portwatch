// Package eventmeta attaches structured metadata key-value pairs to alert events.
package eventmeta

import (
	"fmt"
	"sync"

	"github.com/user/portwatch/internal/alert"
)

// Meta holds arbitrary string metadata for an event keyed by a string.
type Meta map[string]string

// Store manages per-event metadata indexed by event key (port+protocol).
type Store struct {
	mu      sync.RWMutex
	entries map[string]Meta
}

// New returns an initialised Store.
func New() *Store {
	return &Store{entries: make(map[string]Meta)}
}

// eventKey returns a canonical string key for an alert event.
func eventKey(e alert.Event) string {
	return fmt.Sprintf("%s:%d", e.Port.Protocol, e.Port.Number)
}

// Set stores or merges metadata for the event. Existing keys are overwritten.
func (s *Store) Set(e alert.Event, meta Meta) {
	key := eventKey(e)
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.entries[key]; !ok {
		s.entries[key] = make(Meta)
	}
	for k, v := range meta {
		s.entries[key][k] = v
	}
}

// Get retrieves metadata for an event. Returns nil and false when absent.
func (s *Store) Get(e alert.Event) (Meta, bool) {
	key := eventKey(e)
	s.mu.RLock()
	defer s.mu.RUnlock()
	m, ok := s.entries[key]
	if !ok {
		return nil, false
	}
	copy := make(Meta, len(m))
	for k, v := range m {
		copy[k] = v
	}
	return copy, true
}

// Delete removes all metadata for an event.
func (s *Store) Delete(e alert.Event) {
	key := eventKey(e)
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.entries, key)
}

// Reset clears all stored metadata.
func (s *Store) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries = make(map[string]Meta)
}

// Len returns the number of events with metadata stored.
func (s *Store) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.entries)
}
