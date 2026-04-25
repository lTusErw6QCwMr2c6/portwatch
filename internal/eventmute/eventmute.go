// Package eventmute provides a mute registry that suppresses matching events
// for a configurable duration based on port and protocol rules.
package eventmute

import (
	"fmt"
	"sync"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// Rule describes a mute condition.
type Rule struct {
	Name     string
	Port     int
	Protocol string // empty matches any
	Duration time.Duration
}

// entry holds an active mute with its expiry.
type entry struct {
	rule   Rule
	expiry time.Time
}

// Muter suppresses events that match an active mute rule.
type Muter struct {
	mu      sync.RWMutex
	entries map[string]entry
}

// New returns an initialised Muter.
func New() *Muter {
	return &Muter{entries: make(map[string]entry)}
}

// Mute registers a mute rule, silencing matching events for the rule's Duration.
func (m *Muter) Mute(r Rule) error {
	if r.Name == "" {
		return fmt.Errorf("eventmute: rule name must not be empty")
	}
	if r.Duration <= 0 {
		return fmt.Errorf("eventmute: duration must be positive")
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.entries[r.Name] = entry{rule: r, expiry: time.Now().Add(r.Duration)}
	return nil
}

// Unmute removes a named mute rule immediately.
func (m *Muter) Unmute(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.entries, name)
}

// Allow returns true when the event is NOT suppressed by any active rule.
func (m *Muter) Allow(ev alert.Event) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	now := time.Now()
	for _, e := range m.entries {
		if now.After(e.expiry) {
			continue
		}
		if e.rule.Port != ev.Port.Number {
			continue
		}
		if e.rule.Protocol != "" && e.rule.Protocol != ev.Port.Protocol {
			continue
		}
		return false
	}
	return true
}

// Purge removes all expired entries.
func (m *Muter) Purge() {
	m.mu.Lock()
	defer m.mu.Unlock()
	now := time.Now()
	for k, e := range m.entries {
		if now.After(e.expiry) {
			delete(m.entries, k)
		}
	}
}

// Len returns the number of active (non-expired) mute rules.
func (m *Muter) Len() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	now := time.Now()
	count := 0
	for _, e := range m.entries {
		if !now.After(e.expiry) {
			count++
		}
	}
	return count
}
