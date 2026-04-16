package metrics

import (
	"sync"
	"time"
)

// Snapshot holds a point-in-time summary of watcher activity.
type Snapshot struct {
	Cycles      uint64
	Opened      uint64
	Closed      uint64
	Alerts      uint64
	LastCycleAt time.Time
}

// Collector accumulates runtime metrics for the daemon.
type Collector struct {
	mu sync.RWMutex
	s  Snapshot
}

// New returns an initialised Collector.
func New() *Collector {
	return &Collector{}
}

// RecordCycle increments the cycle counter and records the timestamp.
func (c *Collector) RecordCycle() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.s.Cycles++
	c.s.LastCycleAt = time.Now()
}

// RecordOpened adds n newly-opened port events.
func (c *Collector) RecordOpened(n uint64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.s.Opened += n
}

// RecordClosed adds n newly-closed port events.
func (c *Collector) RecordClosed(n uint64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.s.Closed += n
}

// RecordAlert increments the alert counter.
func (c *Collector) RecordAlert() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.s.Alerts++
}

// Read returns a copy of the current snapshot.
func (c *Collector) Read() Snapshot {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.s
}

// Reset zeroes all counters.
func (c *Collector) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.s = Snapshot{}
}
