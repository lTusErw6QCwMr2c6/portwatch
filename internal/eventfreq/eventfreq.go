// Package eventfreq tracks the frequency of events over a sliding time window,
// keyed by port and protocol. It is useful for detecting bursts of activity.
package eventfreq

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// Entry records the count and last-seen time for a key.
type Entry struct {
	Count    int
	LastSeen time.Time
}

// Tracker counts events within a sliding window.
type Tracker struct {
	mu      sync.Mutex
	window  time.Duration
	entries map[string]*Entry
}

// New creates a Tracker with the given sliding window duration.
func New(window time.Duration) *Tracker {
	if window <= 0 {
		window = 30 * time.Second
	}
	return &Tracker{
		window:  window,
		entries: make(map[string]*Entry),
	}
}

func key(e alert.Event) string {
	return e.Port.Protocol + ":" + itoa(e.Port.Number)
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	b := make([]byte, 0, 8)
	for n > 0 {
		b = append([]byte{byte('0' + n%10)}, b...)
		n /= 10
	}
	return string(b)
}

// Record increments the counter for the event's key, evicting stale entries first.
func (t *Tracker) Record(e alert.Event) int {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now()
	k := key(e)

	if ent, ok := t.entries[k]; ok && now.Sub(ent.LastSeen) > t.window {
		delete(t.entries, k)
	}

	ent, ok := t.entries[k]
	if !ok {
		ent = &Entry{}
		t.entries[k] = ent
	}
	ent.Count++
	ent.LastSeen = now
	return ent.Count
}

// Count returns the current count for the event's key, or 0 if absent or expired.
func (t *Tracker) Count(e alert.Event) int {
	t.mu.Lock()
	defer t.mu.Unlock()

	k := key(e)
	ent, ok := t.entries[k]
	if !ok {
		return 0
	}
	if time.Since(ent.LastSeen) > t.window {
		delete(t.entries, k)
		return 0
	}
	return ent.Count
}

// Reset clears all tracked entries.
func (t *Tracker) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.entries = make(map[string]*Entry)
}
