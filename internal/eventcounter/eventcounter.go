package eventcounter

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// Counter tracks event counts per key over a rolling time window.
type Counter struct {
	mu      sync.Mutex
	buckets map[string][]time.Time
	window  time.Duration
}

// New creates a Counter with the given rolling window duration.
func New(window time.Duration) *Counter {
	if window <= 0 {
		window = 60 * time.Second
	}
	return &Counter{
		buckets: make(map[string][]time.Time),
		window:  window,
	}
}

// Record adds an event to the counter for the given key.
func (c *Counter) Record(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now()
	c.buckets[key] = append(c.buckets[key], now)
	c.evict(key, now)
}

// Count returns the number of events recorded for key within the window.
func (c *Counter) Count(key string) int {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.evict(key, time.Now())
	return len(c.buckets[key])
}

// KeyFromEvent returns a string key derived from an alert event.
func KeyFromEvent(e alert.Event) string {
	return e.Port.Protocol + ":" + string(rune(e.Port.Number))
}

// Reset clears all recorded entries.
func (c *Counter) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.buckets = make(map[string][]time.Time)
}

// evict removes timestamps outside the rolling window. Must be called with lock held.
func (c *Counter) evict(key string, now time.Time) {
	cutoff := now.Add(-c.window)
	times := c.buckets[key]
	idx := 0
	for idx < len(times) && times[idx].Before(cutoff) {
		idx++
	}
	if idx > 0 {
		c.buckets[key] = times[idx:]
	}
}
