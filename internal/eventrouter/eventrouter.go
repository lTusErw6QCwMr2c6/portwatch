package eventrouter

import (
	"fmt"
	"sync"

	"github.com/user/portwatch/internal/alert"
)

// Route defines a named handler for matching events.
type Route struct {
	Name    string
	Match   func(e alert.Event) bool
	Handle  func(e alert.Event) error
}

// Router dispatches events to the first matching route.
type Router struct {
	mu     sync.RWMutex
	routes []Route
}

// New returns an empty Router.
func New() *Router {
	return &Router{}
}

// Register adds a route to the router. Routes are evaluated in registration order.
func (r *Router) Register(route Route) error {
	if route.Name == "" {
		return fmt.Errorf("eventrouter: route name must not be empty")
	}
	if route.Match == nil {
		return fmt.Errorf("eventrouter: route %q has nil match function", route.Name)
	}
	if route.Handle == nil {
		return fmt.Errorf("eventrouter: route %q has nil handle function", route.Name)
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.routes = append(r.routes, route)
	return nil
}

// Dispatch sends the event to the first matching route.
// Returns ErrNoMatch if no route matches.
func (r *Router) Dispatch(e alert.Event) error {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, route := range r.routes {
		if route.Match(e) {
			return route.Handle(e)
		}
	}
	return ErrNoMatch
}

// Len returns the number of registered routes.
func (r *Router) Len() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.routes)
}

// ErrNoMatch is returned when no route matches a dispatched event.
var ErrNoMatch = fmt.Errorf("eventrouter: no matching route")
