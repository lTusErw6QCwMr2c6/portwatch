// Package eventexpiry provides a registry that automatically expires events
// after a configurable TTL and invokes a callback when expiry occurs.
package eventexpiry

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// ExpiredFunc is called when an event's TTL elapses.
type ExpiredFunc func(event alert.Event)

// entry holds an event alongside its expiry timer.
type entry struct {
	event alert.Event
	timer *time.Timer
}

// Registry tracks events and fires a callback when they expire.
type Registry struct {
	mu      sync.Mutex
	entries map[string]*entry
	onExpire ExpiredFunc
	ttl      time.Duration
}

// New creates a Registry with the given TTL and expiry callback.
// If onExpire is nil no callback is invoked on expiry.
func New(ttl time.Duration, onExpire ExpiredFunc) *Registry {
	if ttl <= 0 {
		ttl = 30 * time.Second
	}
	return &Registry{
		entries: make(map[string]*entry),
		onExpire: onExpire,
		ttl:      ttl,
	}
}

// Track registers an event. If the key already exists the existing timer is
// reset and the event is replaced.
func (r *Registry) Track(key string, event alert.Event) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if e, ok := r.entries[key]; ok {
		e.timer.Stop()
	}

	timer := time.AfterFunc(r.ttl, func() {
		r.mu.Lock()
		ev, ok := r.entries[key]
		if ok {
			delete(r.entries, key)
		}
		r.mu.Unlock()
		if ok && r.onExpire != nil {
			r.onExpire(ev.event)
		}
	})

	r.entries[key] = &entry{event: event, timer: timer}
}

// Cancel stops the expiry timer for key and removes the event.
// Returns true if the key was present.
func (r *Registry) Cancel(key string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	e, ok := r.entries[key]
	if !ok {
		return false
	}
	e.timer.Stop()
	delete(r.entries, key)
	return true
}

// Len returns the number of currently tracked events.
func (r *Registry) Len() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.entries)
}

// Reset cancels all timers and clears the registry.
func (r *Registry) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, e := range r.entries {
		e.timer.Stop()
	}
	r.entries = make(map[string]*entry)
}
