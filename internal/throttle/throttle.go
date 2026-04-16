package throttle

import (
	"sync"
	"time"
)

// Throttle limits how frequently a keyed action can fire within a time window.
type Throttle struct {
	mu      sync.Mutex
	window  time.Duration
	maxHits int
	hits    map[string][]time.Time
}

// New creates a Throttle that allows up to maxHits occurrences per window.
func New(window time.Duration, maxHits int) *Throttle {
	return &Throttle{
		window:  window,
		maxHits: maxHits,
		hits:    make(map[string][]time.Time),
	}
}

// Allow returns true if the key has not exceeded its hit limit within the window.
func (t *Throttle) Allow(key string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-t.window)

	prev := t.hits[key]
	var recent []time.Time
	for _, ts := range prev {
		if ts.After(cutoff) {
			recent = append(recent, ts)
		}
	}

	if len(recent) >= t.maxHits {
		t.hits[key] = recent
		return false
	}

	t.hits[key] = append(recent, now)
	return true
}

// Reset clears all tracked hits.
func (t *Throttle) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.hits = make(map[string][]time.Time)
}

// Count returns the number of recent hits for the given key.
func (t *Throttle) Count(key string) int {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-t.window)
	count := 0
	for _, ts := range t.hits[key] {
		if ts.After(cutoff) {
			count++
		}
	}
	return count
}
