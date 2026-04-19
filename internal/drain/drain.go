// Package drain provides a graceful event drain that flushes buffered
// events before shutdown, ensuring no events are silently dropped.
package drain

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// Handler is a function that processes a batch of events.
type Handler func(events []alert.Event) error

// Drain buffers events and flushes them on demand or at shutdown.
type Drain struct {
	mu      sync.Mutex
	buf     []alert.Event
	cap     int
	handler Handler
	timeout time.Duration
}

// New creates a new Drain with the given capacity, flush handler, and
// drain timeout used during Stop.
func New(capacity int, timeout time.Duration, h Handler) *Drain {
	if capacity <= 0 {
		capacity = 64
	}
	return &Drain{
		buf:     make([]alert.Event, 0, capacity),
		cap:     capacity,
		handler: h,
		timeout: timeout,
	}
}

// Add enqueues an event. If the buffer is full it triggers an immediate flush.
func (d *Drain) Add(e alert.Event) error {
	d.mu.Lock()
	d.buf = append(d.buf, e)
	full := len(d.buf) >= d.cap
	d.mu.Unlock()
	if full {
		return d.Flush()
	}
	return nil
}

// Flush drains the current buffer by invoking the handler, then clears it.
func (d *Drain) Flush() error {
	d.mu.Lock()
	if len(d.buf) == 0 {
		d.mu.Unlock()
		return nil
	}
	batch := make([]alert.Event, len(d.buf))
	copy(batch, d.buf)
	d.buf = d.buf[:0]
	d.mu.Unlock()
	return d.handler(batch)
}

// Stop flushes remaining events within the configured timeout.
func (d *Drain) Stop() error {
	done := make(chan error, 1)
	go func() { done <- d.Flush() }()
	select {
	case err := <-done:
		return err
	case <-time.After(d.timeout):
		return nil
	}
}

// Len returns the number of buffered events.
func (d *Drain) Len() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return len(d.buf)
}
