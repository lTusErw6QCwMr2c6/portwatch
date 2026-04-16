package history

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// MaxEntries is the default maximum number of history entries to retain.
const MaxEntries = 500

// Entry represents a single recorded alert event with a timestamp.
type Entry struct {
	Timestamp time.Time
	Event     alert.Event
}

// History stores a bounded ring of alert events observed during a run.
type History struct {
	mu      sync.RWMutex
	entries []Entry
	max     int
}

// New creates a History with the given capacity cap.
// If cap is <= 0 the default MaxEntries is used.
func New(cap int) *History {
	if cap <= 0 {
		cap = MaxEntries
	}
	return &History{
		entries: make([]Entry, 0, cap),
		max:     cap,
	}
}

// Add appends an event to the history. If the buffer is full the oldest
// entry is evicted.
func (h *History) Add(e alert.Event) {
	h.mu.Lock()
	defer h.mu.Unlock()
	entry := Entry{Timestamp: time.Now(), Event: e}
	if len(h.entries) >= h.max {
		// evict oldest
		h.entries = append(h.entries[1:], entry)
	} else {
		h.entries = append(h.entries, entry)
	}
}

// All returns a copy of all stored entries, oldest first.
func (h *History) All() []Entry {
	h.mu.RLock()
	defer h.mu.RUnlock()
	copy := make([]Entry, len(h.entries))
	for i, e := range h.entries {
		copy[i] = e
	}
	return copy
}

// Len returns the current number of stored entries.
func (h *History) Len() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.entries)
}

// Clear removes all stored entries.
func (h *History) Clear() {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.entries = h.entries[:0]
}
