// Package suppress provides a mechanism to temporarily suppress
// repeated alerts for known or acknowledged port events.
package suppress

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// Entry holds suppression metadata for a key.
type Entry struct {
	Until  time.Time
	Reason string
}

// Suppressor stores suppression rules keyed by "proto:port".
type Suppressor struct {
	mu      sync.RWMutex
	entries map[string]Entry
	now     func() time.Time
}

// New returns a new Suppressor.
func New() *Suppressor {
	return &Suppressor{
		entries: make(map[string]Entry),
		now:     time.Now,
	}
}

// Suppress marks the given key as suppressed until the given duration elapses.
func (s *Suppressor) Suppress(key string, duration time.Duration, reason string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries[key] = Entry{
		Until:  s.now().Add(duration),
		Reason: reason,
	}
}

// IsSuppressed returns true if the key is currently suppressed.
func (s *Suppressor) IsSuppressed(key string) (bool, string) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.entries[key]
	if !ok {
		return false, ""
	}
	if s.now().Before(e.Until) {
		return true, e.Reason
	}
	return false, ""
}

// Allow returns true if the event should be forwarded (not suppressed).
func (s *Suppressor) Allow(ev alert.Event) bool {
	key := ev.Port.Protocol + ":" + ev.Port.Address
	suppressed, _ := s.IsSuppressed(key)
	return !suppressed
}

// Remove lifts suppression for the given key immediately.
func (s *Suppressor) Remove(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.entries, key)
}

// Reset clears all suppression entries.
func (s *Suppressor) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries = make(map[string]Entry)
}
