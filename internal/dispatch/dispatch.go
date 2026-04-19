package dispatch

import (
	"fmt"
	"sync"

	"github.com/user/portwatch/internal/alert"
)

// Handler is a function that processes an alert event.
type Handler func(event alert.Event) error

// Router dispatches alert events to registered named handlers.
type Router struct {
	mu       sync.RWMutex
	handlers map[string]Handler
}

// New returns an empty Router.
func New() *Router {
	return &Router{
		handlers: make(map[string]Handler),
	}
}

// Register adds a named handler to the router.
func (r *Router) Register(name string, h Handler) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.handlers[name] = h
}

// Deregister removes a handler by name.
func (r *Router) Deregister(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.handlers, name)
}

// Dispatch sends the event to all registered handlers and collects errors.
func (r *Router) Dispatch(event alert.Event) []error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var errs []error
	for name, h := range r.handlers {
		if err := h(event); err != nil {
			errs = append(errs, fmt.Errorf("handler %q: %w", name, err))
		}
	}
	return errs
}

// Len returns the number of registered handlers.
func (r *Router) Len() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.handlers)
}
