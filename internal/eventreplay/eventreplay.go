package eventreplay

import (
	"errors"
	"sync"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// ErrEmptyBuffer is returned when the replay buffer contains no events.
var ErrEmptyBuffer = errors.New("eventreplay: buffer is empty")

// Entry wraps an alert event with the time it was recorded.
type Entry struct {
	Event     alert.Event
	RecordedAt time.Time
}

// Buffer is a bounded, in-memory replay buffer for alert events.
type Buffer struct {
	mu      sync.RWMutex
	entries []Entry
	cap     int
}

// New returns a Buffer with the given maximum capacity.
// If cap is <= 0 it defaults to 256.
func New(cap int) *Buffer {
	if cap <= 0 {
		cap = 256
	}
	return &Buffer{cap: cap}
}

// Add appends an event to the buffer, evicting the oldest entry when full.
func (b *Buffer) Add(e alert.Event) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if len(b.entries) >= b.cap {
		b.entries = b.entries[1:]
	}
	b.entries = append(b.entries, Entry{Event: e, RecordedAt: time.Now()})
}

// Since returns all entries recorded at or after t.
func (b *Buffer) Since(t time.Time) []Entry {
	b.mu.RLock()
	defer b.mu.RUnlock()
	var out []Entry
	for _, e := range b.entries {
		if !e.RecordedAt.Before(t) {
			out = append(out, e)
		}
	}
	return out
}

// All returns a copy of every entry in the buffer.
func (b *Buffer) All() []Entry {
	b.mu.RLock()
	defer b.mu.RUnlock()
	out := make([]Entry, len(b.entries))
	copy(out, b.entries)
	return out
}

// Len returns the current number of buffered entries.
func (b *Buffer) Len() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.entries)
}

// Reset clears all buffered entries.
func (b *Buffer) Reset() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.entries = b.entries[:0]
}
