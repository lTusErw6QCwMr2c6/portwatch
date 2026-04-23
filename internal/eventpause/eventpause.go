// Package eventpause provides a pauseable gate for event processing.
// When paused, events submitted via Allow are held and the caller is
// notified that the gate is closed. Callers can resume processing by
// calling Resume, which re-opens the gate.
package eventpause

import (
	"errors"
	"sync"
)

// ErrPaused is returned by Allow when the gate is currently paused.
var ErrPaused = errors.New("eventpause: gate is paused")

// Gate is a pauseable gate that controls whether events are allowed through.
type Gate struct {
	mu     sync.RWMutex
	paused bool
	count  int64 // total events allowed through
	drops  int64 // total events dropped while paused
}

// New returns a new Gate in the resumed (open) state.
func New() *Gate {
	return &Gate{}
}

// Pause closes the gate. Subsequent calls to Allow will return ErrPaused.
func (g *Gate) Pause() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.paused = true
}

// Resume opens the gate. Subsequent calls to Allow will succeed.
func (g *Gate) Resume() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.paused = false
}

// IsPaused reports whether the gate is currently paused.
func (g *Gate) IsPaused() bool {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.paused
}

// Allow checks whether the gate is open. It returns nil when the gate is
// open and increments the allowed counter. It returns ErrPaused when the
// gate is closed and increments the drop counter.
func (g *Gate) Allow() error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.paused {
		g.drops++
		return ErrPaused
	}
	g.count++
	return nil
}

// Stats returns the total allowed and dropped event counts.
func (g *Gate) Stats() (allowed, dropped int64) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.count, g.drops
}

// Reset clears the allowed and dropped counters and opens the gate.
func (g *Gate) Reset() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.paused = false
	g.count = 0
	g.drops = 0
}
