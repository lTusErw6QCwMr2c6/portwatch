// Package circuitbreaker provides a simple circuit breaker to halt
// repeated failing operations until a cooldown period has elapsed.
package circuitbreaker

import (
	"errors"
	"sync"
	"time"
)

// ErrOpen is returned when the circuit is open and calls are rejected.
var ErrOpen = errors.New("circuit breaker is open")

// State represents the current circuit breaker state.
type State int

const (
	StateClosed State = iota
	StateOpen
	StateHalfOpen
)

// Breaker tracks consecutive failures and opens the circuit when the
// failure threshold is exceeded.
type Breaker struct {
	mu           sync.Mutex
	threshold    int
	cooldown     time.Duration
	failures     int
	state        State
	openedAt     time.Time
}

// New creates a Breaker that opens after threshold consecutive failures
// and attempts recovery after cooldown.
func New(threshold int, cooldown time.Duration) *Breaker {
	if threshold <= 0 {
		threshold = 3
	}
	if cooldown <= 0 {
		cooldown = 30 * time.Second
	}
	return &Breaker{
		threshold: threshold,
		cooldown:  cooldown,
		state:     StateClosed,
	}
}

// Allow returns nil if the call should proceed, or ErrOpen if the circuit
// is open. It transitions to half-open once the cooldown expires.
func (b *Breaker) Allow() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	switch b.state {
	case StateOpen:
		if time.Since(b.openedAt) >= b.cooldown {
			b.state = StateHalfOpen
			return nil
		}
		return ErrOpen
	default:
		return nil
	}
}

// RecordSuccess resets the breaker to closed state.
func (b *Breaker) RecordSuccess() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures = 0
	b.state = StateClosed
}

// RecordFailure increments the failure count and opens the circuit if the
// threshold is reached.
func (b *Breaker) RecordFailure() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures++
	if b.failures >= b.threshold {
		b.state = StateOpen
		b.openedAt = time.Now()
	}
}

// State returns the current state of the breaker.
func (b *Breaker) State() State {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.state
}

// Reset forces the breaker back to closed state and clears failures.
func (b *Breaker) Reset() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures = 0
	b.state = StateClosed
}
