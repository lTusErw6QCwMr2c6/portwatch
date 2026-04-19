// Package sequence tracks the order of port events to detect rapid sequential scanning.
package sequence

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// Entry holds a recorded event with its timestamp.
type Entry struct {
	Event     alert.Event
	RecordedAt time.Time
}

// Tracker records events in order and can detect sequential port sweeps.
type Tracker struct {
	mu      sync.Mutex
	entries []Entry
	window  time.Duration
	maxSize int
}

// New returns a Tracker that retains events within the given window.
func New(window time.Duration, maxSize int) *Tracker {
	if maxSize <= 0 {
		maxSize = 256
	}
	return &Tracker{window: window, maxSize: maxSize}
}

// Add records an event.
func (t *Tracker) Add(e alert.Event) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.evict()
	if len(t.entries) >= t.maxSize {
		t.entries = t.entries[1:]
	}
	t.entries = append(t.entries, Entry{Event: e, RecordedAt: time.Now()})
}

// Recent returns all events still within the tracking window.
func (t *Tracker) Recent() []Entry {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.evict()
	out := make([]Entry, len(t.entries))
	copy(out, t.entries)
	return out
}

// Count returns the number of events currently in the window.
func (t *Tracker) Count() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.evict()
	return len(t.entries)
}

// Reset clears all recorded entries.
func (t *Tracker) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.entries = nil
}

func (t *Tracker) evict() {
	cutoff := time.Now().Add(-t.window)
	i := 0
	for i < len(t.entries) && t.entries[i].RecordedAt.Before(cutoff) {
		i++
	}
	t.entries = t.entries[i:]
}
