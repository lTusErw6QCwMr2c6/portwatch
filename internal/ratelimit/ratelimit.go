package ratelimit

import (
	"sync"
	"time"
)

// Limiter suppresses repeated alerts for the same port within a cooldown window.
type Limiter struct {
	mu       sync.Mutex
	cooldown time.Duration
	seen     map[string]time.Time
}

// New creates a Limiter with the given cooldown duration.
// Events for the same key will be suppressed until the cooldown expires.
func New(cooldown time.Duration) *Limiter {
	return &Limiter{
		cooldown: cooldown,
		seen:     make(map[string]time.Time),
	}
}

// Allow returns true if the event identified by key should be allowed through.
// It returns false if the same key was seen within the cooldown window.
func (l *Limiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	if last, ok := l.seen[key]; ok {
		if now.Sub(last) < l.cooldown {
			return false
		}
	}
	l.seen[key] = now
	return true
}

// Reset clears all tracked keys, allowing all events through again.
func (l *Limiter) Reset() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.seen = make(map[string]time.Time)
}

// Prune removes entries whose cooldown has already expired,
// keeping memory usage bounded during long-running operation.
func (l *Limiter) Prune() {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	for key, last := range l.seen {
		if now.Sub(last) >= l.cooldown {
			delete(l.seen, key)
		}
	}
}
