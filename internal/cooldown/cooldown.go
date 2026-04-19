// Package cooldown provides per-key cooldown enforcement to prevent
// repeated actions from firing too frequently.
package cooldown

import (
	"sync"
	"time"
)

// Cooldown tracks the last activation time per key and enforces a minimum
// duration between successive activations.
type Cooldown struct {
	mu       sync.Mutex
	last     map[string]time.Time
	duration time.Duration
}

// New creates a Cooldown with the given duration between allowed activations.
func New(d time.Duration) *Cooldown {
	return &Cooldown{
		last:     make(map[string]time.Time),
		duration: d,
	}
}

// Allow returns true and records the activation time if the key has not been
// activated within the cooldown window. Returns false otherwise.
func (c *Cooldown) Allow(key string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now()
	if t, ok := c.last[key]; ok && now.Sub(t) < c.duration {
		return false
	}
	c.last[key] = now
	return true
}

// Remaining returns how much cooldown time is left for the given key.
// Returns zero if the key is not in cooldown.
func (c *Cooldown) Remaining(key string) time.Duration {
	c.mu.Lock()
	defer c.mu.Unlock()
	t, ok := c.last[key]
	if !ok {
		return 0
	}
	remaining := c.duration - time.Since(t)
	if remaining < 0 {
		return 0
	}
	return remaining
}

// Reset clears all recorded activation times.
func (c *Cooldown) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.last = make(map[string]time.Time)
}
