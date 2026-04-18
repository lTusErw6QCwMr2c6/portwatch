// Package dedupe suppresses duplicate alert events within a sliding window.
package dedupe

import (
	"crypto/sha256"
	"fmt"
	"sync"
	"time"
)

// Filter tracks recently seen event fingerprints and drops duplicates.
type Filter struct {
	mu      sync.Mutex
	seen    map[string]time.Time
	window  time.Duration
	now     func() time.Time
}

// New returns a Filter that suppresses repeated events within window.
func New(window time.Duration) *Filter {
	return &Filter{
		seen:   make(map[string]time.Time),
		window: window,
		now:    time.Now,
	}
}

// Allow returns true if the event identified by key has not been seen
// within the configured window. Allowed events are recorded.
func (f *Filter) Allow(key string) bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.evict()
	fp := fingerprint(key)
	if _, exists := f.seen[fp]; exists {
		return false
	}
	f.seen[fp] = f.now()
	return true
}

// Reset clears all recorded fingerprints.
func (f *Filter) Reset() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.seen = make(map[string]time.Time)
}

// Len returns the number of active fingerprints.
func (f *Filter) Len() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.evict()
	return len(f.seen)
}

// evict removes entries older than the window. Caller must hold mu.
func (f *Filter) evict() {
	cutoff := f.now().Add(-f.window)
	for k, t := range f.seen {
		if t.Before(cutoff) {
			delete(f.seen, k)
		}
	}
}

func fingerprint(key string) string {
	h := sha256.Sum256([]byte(key))
	return fmt.Sprintf("%x", h[:8])
}
