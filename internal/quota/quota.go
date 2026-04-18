// Package quota enforces per-key event volume limits over a sliding window.
package quota

import (
	"sync"
	"time"
)

// Config holds quota settings.
type Config struct {
	Limit  int
	Window time.Duration
}

// DefaultConfig returns sensible defaults.
func DefaultConfig() Config {
	return Config{
		Limit:  100,
		Window: time.Minute,
	}
}

type entry struct {
	count     int
	windowEnd time.Time
}

// Quota tracks event counts per key and enforces a limit per window.
type Quota struct {
	mu      sync.Mutex
	cfg     Config
	buckets map[string]*entry
}

// New creates a new Quota enforcer.
func New(cfg Config) *Quota {
	return &Quota{
		cfg:     cfg,
		buckets: make(map[string]*entry),
	}
}

// Allow returns true if the key is within quota, false if exceeded.
func (q *Quota) Allow(key string) bool {
	q.mu.Lock()
	defer q.mu.Unlock()

	now := time.Now()
	e, ok := q.buckets[key]
	if !ok || now.After(e.windowEnd) {
		q.buckets[key] = &entry{count: 1, windowEnd: now.Add(q.cfg.Window)}
		return true
	}
	if e.count >= q.cfg.Limit {
		return false
	}
	e.count++
	return true
}

// Remaining returns how many events the key may still emit in the current window.
func (q *Quota) Remaining(key string) int {
	q.mu.Lock()
	defer q.mu.Unlock()

	now := time.Now()
	e, ok := q.buckets[key]
	if !ok || now.After(e.windowEnd) {
		return q.cfg.Limit
	}
	r := q.cfg.Limit - e.count
	if r < 0 {
		return 0
	}
	return r
}

// Reset clears all quota state.
func (q *Quota) Reset() {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.buckets = make(map[string]*entry)
}
