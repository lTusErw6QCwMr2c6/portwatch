// Package graceful provides utilities for orderly shutdown of portwatch components.
package graceful

import (
	"context"
	"sync"
	"time"
)

// ShutdownFunc is a function called during shutdown.
type ShutdownFunc func(ctx context.Context) error

// Shutdown coordinates an ordered shutdown sequence, calling registered
// handlers in LIFO order with a shared deadline context.
type Shutdown struct {
	mu       sync.Mutex
	handlers []namedHandler
	timeout  time.Duration
}

type namedHandler struct {
	name string
	fn   ShutdownFunc
}

// New returns a Shutdown with the given timeout applied to each handler.
func New(timeout time.Duration) *Shutdown {
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	return &Shutdown{timeout: timeout}
}

// Register adds a shutdown handler with an identifying name.
// Handlers are called in reverse registration order.
func (s *Shutdown) Register(name string, fn ShutdownFunc) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.handlers = append(s.handlers, namedHandler{name: name, fn: fn})
}

// Result holds the outcome of a single shutdown handler.
type Result struct {
	Name string
	Err  error
}

// Run executes all registered handlers in LIFO order.
// Each handler receives a fresh context bounded by the configured timeout.
// All handlers are attempted regardless of prior errors.
func (s *Shutdown) Run() []Result {
	s.mu.Lock()
	handlers := make([]namedHandler, len(s.handlers))
	copy(handlers, s.handlers)
	s.mu.Unlock()

	results := make([]Result, 0, len(handlers))

	for i := len(handlers) - 1; i >= 0; i-- {
		h := handlers[i]
		ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
		err := h.fn(ctx)
		cancel()
		results = append(results, Result{Name: h.name, Err: err})
	}

	return results
}

// HasErrors returns true if any result contains a non-nil error.
func HasErrors(results []Result) bool {
	for _, r := range results {
		if r.Err != nil {
			return true
		}
	}
	return false
}
