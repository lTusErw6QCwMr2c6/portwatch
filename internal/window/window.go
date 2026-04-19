// Package window provides a sliding time-window counter for tracking
// event frequency over a rolling duration.
package window

import (
	"sync"
	"time"
)

// Counter tracks event occurrences within a sliding time window per key.
type Counter struct {
	mu       sync.Mutex
	window   time.Duration
	buckets  map[string][]time.Time
}

// New returns a Counter with the given sliding window duration.
func New(window time.Duration) *Counter {
	return &Counter{
		window:  window,
		buckets: make(map[string][]time.Time),
	}
}

// Add records an event for the given key at the current time.
func (c *Counter) Add(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now()
	c.buckets[key] = append(c.evict(key, now), now)
}

// Count returns the number of events recorded for key within the window.
func (c *Counter) Count(key string) int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.evict(key, time.Now()))
}

// Reset clears all recorded events.
func (c *Counter) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.buckets = make(map[string][]time.Time)
}

// evict removes timestamps outside the window and updates the bucket.
// Must be called with the lock held.
func (c *Counter) evict(key string, now time.Time) []time.Time {
	cutoff := now.Add(-c.window)
	times := c.buckets[key]
	i := 0
	for i < len(times) && times[i].Before(cutoff) {
		i++
	}
	valid := times[i:]
	c.buckets[key] = valid
	return valid
}
