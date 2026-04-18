// Package buffer provides a fixed-size ring buffer for storing recent port events.
package buffer

import (
	"sync"

	"github.com/user/portwatch/internal/alert"
)

// Buffer is a thread-safe ring buffer of alert events.
type Buffer struct {
	mu       sync.Mutex
	items    []alert.Event
	cap      int
	head     int
	count    int
}

// New creates a new Buffer with the given capacity.
func New(capacity int) *Buffer {
	if capacity <= 0 {
		capacity = 64
	}
	return &Buffer{
		items: make([]alert.Event, capacity),
		cap:   capacity,
	}
}

// Add inserts an event into the ring buffer, overwriting the oldest entry when full.
func (b *Buffer) Add(e alert.Event) {
	b.mu.Lock()
	defer b.mu.Unlock()
	index := (b.head + b.count) % b.cap
	if b.count == b.cap {
		// overwrite oldest
		b.items[b.head] = e
		b.head = (b.head + 1) % b.cap
	} else {
		b.items[index] = e
		b.count++
	}
}

// All returns a snapshot of all buffered events in insertion order.
func (b *Buffer) All() []alert.Event {
	b.mu.Lock()
	defer b.mu.Unlock()
	out := make([]alert.Event, b.count)
	for i := 0; i < b.count; i++ {
		out[i] = b.items[(b.head+i)%b.cap]
	}
	return out
}

// Len returns the number of events currently in the buffer.
func (b *Buffer) Len() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.count
}

// Reset clears all events from the buffer.
func (b *Buffer) Reset() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.head = 0
	b.count = 0
}
