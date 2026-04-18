package correlation

import (
	"fmt"
	"sync"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// Rule defines a pattern that triggers a correlated alert.
type Rule struct {
	Name     string
	MinPorts int
	Window   time.Duration
}

// Match represents a correlated alert match.
type Match struct {
	Rule      string
	Events    []alert.Event
	DetectedAt time.Time
}

func (m Match) String() string {
	return fmt.Sprintf("[CORRELATION] rule=%s ports=%d at=%s",
		m.Rule, len(m.Events), m.DetectedAt.Format(time.RFC3339))
}

// Correlator groups events within a time window and matches rules.
type Correlator struct {
	mu     sync.Mutex
	rules  []Rule
	bucket []alert.Event
}

// New returns a Correlator with the given rules.
func New(rules []Rule) *Correlator {
	return &Correlator{rules: rules}
}

// Add appends an event and evicts entries outside the longest window.
func (c *Correlator) Add(e alert.Event) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.bucket = append(c.bucket, e)
	c.evict()
}

// Evaluate checks all rules against the current bucket.
func (c *Correlator) Evaluate() []Match {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.evict()
	now := time.Now()
	var matches []Match
	for _, r := range c.rules {
		cutoff := now.Add(-r.Window)
		var window []alert.Event
		for _, ev := range c.bucket {
			if ev.Time.After(cutoff) {
				window = append(window, ev)
			}
		}
		if len(window) >= r.MinPorts {
			matches = append(matches, Match{Rule: r.Name, Events: window, DetectedAt: now})
		}
	}
	return matches
}

// Reset clears the internal event bucket.
func (c *Correlator) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.bucket = nil
}

func (c *Correlator) evict() {
	if len(c.rules) == 0 {
		return
	}
	var longest time.Duration
	for _, r := range c.rules {
		if r.Window > longest {
			longest = r.Window
		}
	}
	cutoff := time.Now().Add(-longest)
	filtered := c.bucket[:0]
	for _, e := range c.bucket {
		if e.Time.After(cutoff) {
			filtered = append(filtered, e)
		}
	}
	c.bucket = filtered
}
