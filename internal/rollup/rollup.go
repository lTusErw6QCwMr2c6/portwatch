// Package rollup groups rapid bursts of port events into a single summary.
package rollup

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// Group holds aggregated events within a window.
type Group struct {
	Opened []alert.Event
	Closed []alert.Event
	WindowEnd time.Time
}

// Rollup buffers events and flushes them as groups after a quiet window.
type Rollup struct {
	mu      sync.Mutex
	window  time.Duration
	buf     []alert.Event
	timer   *time.Timer
	onFlush func(Group)
}

// New creates a Rollup that calls onFlush after window of inactivity.
func New(window time.Duration, onFlush func(Group)) *Rollup {
	return &Rollup{
		window:  window,
		onFlush: onFlush,
	}
}

// Add buffers an event and resets the flush timer.
func (r *Rollup) Add(e alert.Event) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.buf = append(r.buf, e)
	if r.timer != nil {
		r.timer.Reset(r.window)
		return
	}
	r.timer = time.AfterFunc(r.window, r.flush)
}

func (r *Rollup) flush() {
	r.mu.Lock()
	events := r.buf
	r.buf = nil
	r.timer = nil
	r.mu.Unlock()

	g := Group{WindowEnd: time.Now()}
	for _, e := range events {
		switch e.Kind {
		case alert.Opened:
			g.Opened = append(g.Opened, e)
		case alert.Closed:
			g.Closed = append(g.Closed, e)
		}
	}
	r.onFlush(g)
}

// Flush forces an immediate flush regardless of the window.
func (r *Rollup) Flush() {
	r.mu.Lock()
	if r.timer != nil {
		r.timer.Stop()
		r.timer = nil
	}
	r.mu.Unlock()
	r.flush()
}

// Len returns the number of buffered events.
func (r *Rollup) Len() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.buf)
}
