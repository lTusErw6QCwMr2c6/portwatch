// Package eventbridge routes events between named sources and sinks.
package eventbridge

import (
	"errors"
	"fmt"
	"sync"

	"github.com/user/portwatch/internal/alert"
)

// Handler is a function that receives an alert event.
type Handler func(event alert.Event)

// Bridge routes events from registered sources to bound sinks.
type Bridge struct {
	mu       sync.RWMutex
	sinks    map[string][]Handler
	sources  map[string]struct{}
}

// New returns an initialised Bridge.
func New() *Bridge {
	return &Bridge{
		sinks:   make(map[string][]Handler),
		sources: make(map[string]struct{}),
	}
}

// RegisterSource declares a named event source.
func (b *Bridge) RegisterSource(name string) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if _, ok := b.sources[name]; ok {
		return fmt.Errorf("eventbridge: source %q already registered", name)
	}
	b.sources[name] = struct{}{}
	return nil
}

// Subscribe attaches a handler to events emitted by source.
func (b *Bridge) Subscribe(source string, h Handler) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if _, ok := b.sources[source]; !ok {
		return fmt.Errorf("eventbridge: unknown source %q", source)
	}
	b.sinks[source] = append(b.sinks[source], h)
	return nil
}

// Emit sends an event to all handlers subscribed to source.
func (b *Bridge) Emit(source string, event alert.Event) error {
	b.mu.RLock()
	defer b.mu.RUnlock()
	if _, ok := b.sources[source]; !ok {
		return errors.New("eventbridge: emit on unknown source " + source)
	}
	for _, h := range b.sinks[source] {
		h(event)
	}
	return nil
}

// Sources returns the names of all registered sources.
func (b *Bridge) Sources() []string {
	b.mu.RLock()
	defer b.mu.RUnlock()
	out := make([]string, 0, len(b.sources))
	for k := range b.sources {
		out = append(out, k)
	}
	return out
}
