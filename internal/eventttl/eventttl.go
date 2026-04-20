// Package eventttl provides time-to-live expiry tracking for events.
// Events that exceed their TTL are considered expired and can be evicted.
package eventttl

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// Entry holds an event alongside its expiry deadline.
type Entry struct {
	Event    alert.Event
	ExpiresAt time.Time
}

// Expired reports whether the entry has passed its deadline.
func (e Entry) Expired() bool {
	return time.Now().After(e.ExpiresAt)
}

// Store tracks events with individual TTLs.
type Store struct {
	mu      sync.Mutex
	entries map[string]Entry
	ttl     time.Duration
}

// New returns a Store that assigns ttl to every added event.
func New(ttl time.Duration) *Store {
	return &Store{
		entries: make(map[string]Entry),
		ttl:     ttl,
	}
}

// Add inserts or refreshes an event using its fingerprint as key.
func (s *Store) Add(key string, ev alert.Event) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries[key] = Entry{
		Event:    ev,
		ExpiresAt: time.Now().Add(s.ttl),
	}
}

// Get returns the entry for key and whether it was found and is still live.
func (s *Store) Get(key string) (Entry, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.entries[key]
	if !ok || e.Expired() {
		return Entry{}, false
	}
	return e, true
}

// Evict removes all entries whose TTL has elapsed and returns them.
func (s *Store) Evict() []Entry {
	s.mu.Lock()
	defer s.mu.Unlock()
	var expired []Entry
	for k, e := range s.entries {
		if e.Expired() {
			expired = append(expired, e)
			delete(s.entries, k)
		}
	}
	return expired
}

// Len returns the number of live (non-expired) entries currently held.
func (s *Store) Len() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	count := 0
	for _, e := range s.entries {
		if !e.Expired() {
			count++
		}
	}
	return count
}

// Reset removes all entries regardless of expiry.
func (s *Store) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries = make(map[string]Entry)
}
