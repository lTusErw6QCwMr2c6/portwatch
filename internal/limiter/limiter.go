// Package limiter enforces a maximum number of active monitored ports at once.
package limiter

import (
	"errors"
	"sync"

	"github.com/user/portwatch/internal/alert"
)

// ErrLimitExceeded is returned when the active port count would exceed the cap.
var ErrLimitExceeded = errors.New("limiter: active port limit exceeded")

// Limiter tracks how many ports are currently open and enforces a hard cap.
type Limiter struct {
	mu    sync.Mutex
	max   int
	count int
}

// New creates a Limiter with the given maximum active port count.
func New(max int) *Limiter {
	if max <= 0 {
		max = 1024
	}
	return &Limiter{max: max}
}

// Apply processes a slice of alert events, updating the active count and
// returning only the events that are within the configured limit.
// Opened events that would breach the cap are dropped with ErrLimitExceeded.
func (l *Limiter) Apply(events []alert.Event) ([]alert.Event, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	var allowed []alert.Event
	var firstErr error

	for _, e := range events {
		switch e.Type {
		case alert.Opened:
			if l.count >= l.max {
				if firstErr == nil {
					firstErr = ErrLimitExceeded
				}
				continue
			}
			l.count++
			allowed = append(allowed, e)
		case alert.Closed:
			if l.count > 0 {
				l.count--
			}
			allowed = append(allowed, e)
		default:
			allowed = append(allowed, e)
		}
	}
	return allowed, firstErr
}

// Count returns the current number of tracked open ports.
func (l *Limiter) Count() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.count
}

// Reset clears the active count.
func (l *Limiter) Reset() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.count = 0
}
