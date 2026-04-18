// Package trend tracks port event frequency over a sliding window and
// exposes a simple rate (events per minute) for each port key.
package trend

import (
	"sync"
	"time"
)

// Entry holds timestamped event counts for a single key.
type Entry struct {
	times []time.Time
}

// Tracker records events and computes per-key rates over a sliding window.
type Tracker struct {
	mu     sync.Mutex
	window time.Duration
	bucket map[string]*Entry
}

// New returns a Tracker with the given sliding window duration.
func New(window time.Duration) *Tracker {
	if window <= 0 {
		window = time.Minute
	}
	return &Tracker{
		window: window,
		bucket: make(map[string]*Entry),
	}
}

// Record registers one event for key at the current time.
func (t *Tracker) Record(key string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	e, ok := t.bucket[key]
	if !ok {
		e = &Entry{}
		t.bucket[key] = e
	}
	now := time.Now()
	e.times = append(e.times, now)
	t.evict(e, now)
}

// Rate returns the number of events recorded for key within the window.
func (t *Tracker) Rate(key string) int {
	t.mu.Lock()
	defer t.mu.Unlock()
	e, ok := t.bucket[key]
	if !ok {
		return 0
	}
	t.evict(e, time.Now())
	return len(e.times)
}

// Reset clears all recorded events.
func (t *Tracker) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.bucket = make(map[string]*Entry)
}

// evict removes timestamps outside the sliding window. Must be called with lock held.
func (t *Tracker) evict(e *Entry, now time.Time) {
	cutoff := now.Add(-t.window)
	i := 0
	for i < len(e.times) && e.times[i].Before(cutoff) {
		i++
	}
	e.times = e.times[i:]
}
