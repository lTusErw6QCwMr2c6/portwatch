package eventlog

import (
	"fmt"
	"sync"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// Entry represents a single logged event with metadata.
type Entry struct {
	ID        string
	Timestamp time.Time
	Event     alert.Event
	Source    string
	Tags      []string
}

func (e Entry) String() string {
	return fmt.Sprintf("[%s] %s %s (source=%s)",
		e.Timestamp.Format(time.RFC3339),
		e.ID,
		e.Event.String(),
		e.Source,
	)
}

// Log is an in-memory ordered event log with a configurable capacity.
type Log struct {
	mu       sync.RWMutex
	entries  []Entry
	capacity int
}

// New creates a new Log with the given capacity.
func New(capacity int) *Log {
	if capacity <= 0 {
		capacity = 256
	}
	return &Log{capacity: capacity}
}

// Append adds an entry to the log, evicting the oldest if at capacity.
func (l *Log) Append(e Entry) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if len(l.entries) >= l.capacity {
		l.entries = l.entries[1:]
	}
	l.entries = append(l.entries, e)
}

// All returns a copy of all entries in insertion order.
func (l *Log) All() []Entry {
	l.mu.RLock()
	defer l.mu.RUnlock()
	out := make([]Entry, len(l.entries))
	copy(out, l.entries)
	return out
}

// Since returns entries with a timestamp after t.
func (l *Log) Since(t time.Time) []Entry {
	l.mu.RLock()
	defer l.mu.RUnlock()
	var out []Entry
	for _, e := range l.entries {
		if e.Timestamp.After(t) {
			out = append(out, e)
		}
	}
	return out
}

// Len returns the current number of entries.
func (l *Log) Len() int {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return len(l.entries)
}

// Clear removes all entries.
func (l *Log) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.entries = nil
}
