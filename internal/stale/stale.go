// Package stale detects ports that have been open beyond a configured duration.
package stale

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// Entry tracks when a port was first seen open.
type Entry struct {
	FirstSeen time.Time
	Port      int
	Protocol  string
}

// Detector tracks open ports and flags those exceeding a max age.
type Detector struct {
	mu      sync.Mutex
	entries map[string]Entry
	maxAge  time.Duration
}

// New creates a Detector with the given max age threshold.
func New(maxAge time.Duration) *Detector {
	return &Detector{
		entries: make(map[string]Entry),
		maxAge:  maxAge,
	}
}

func key(e alert.Event) string {
	return e.Protocol + ":" + string(rune(e.Port))
}

// Track updates internal state based on an event.
// Opened events register a port; Closed events remove it.
func (d *Detector) Track(e alert.Event) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if e.Type == alert.Opened {
		k := key(e)
		if _, exists := d.entries[k]; !exists {
			d.entries[k] = Entry{
				FirstSeen: time.Now(),
				Port:      e.Port,
				Protocol:  e.Protocol,
			}
		}
	} else {
		delete(d.entries, key(e))
	}
}

// Stale returns all entries that have been open longer than maxAge.
func (d *Detector) Stale() []Entry {
	d.mu.Lock()
	defer d.mu.Unlock()
	now := time.Now()
	var out []Entry
	for _, e := range d.entries {
		if now.Sub(e.FirstSeen) >= d.maxAge {
			out = append(out, e)
		}
	}
	return out
}

// Reset clears all tracked entries.
func (d *Detector) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.entries = make(map[string]Entry)
}
