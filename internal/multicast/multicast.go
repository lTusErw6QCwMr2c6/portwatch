// Package multicast fans out events to multiple named subscribers.
package multicast

import (
	"fmt"
	"sync"

	"github.com/user/portwatch/internal/alert"
)

// Subscriber receives events from the bus.
type Subscriber struct {
	Name string
	Ch   chan alert.Event
}

// Bus distributes events to all registered subscribers.
type Bus struct {
	mu   sync.RWMutex
	subs map[string]*Subscriber
	buf  int
}

// New creates a Bus with the given channel buffer size.
func New(buf int) *Bus {
	if buf < 1 {
		buf = 1
	}
	return &Bus{subs: make(map[string]*Subscriber), buf: buf}
}

// Subscribe registers a named subscriber and returns its channel.
func (b *Bus) Subscribe(name string) (<-chan alert.Event, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if _, exists := b.subs[name]; exists {
		return nil, fmt.Errorf("multicast: subscriber %q already registered", name)
	}
	s := &Subscriber{Name: name, Ch: make(chan alert.Event, b.buf)}
	b.subs[name] = s
	return s.Ch, nil
}

// Unsubscribe removes a subscriber and closes its channel.
func (b *Bus) Unsubscribe(name string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if s, ok := b.subs[name]; ok {
		close(s.Ch)
		delete(b.subs, name)
	}
}

// Publish sends an event to all subscribers, dropping if their buffer is full.
func (b *Bus) Publish(ev alert.Event) int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	delivered := 0
	for _, s := range b.subs {
		select {
		case s.Ch <- ev:
			delivered++
		default:
		}
	}
	return delivered
}

// Len returns the number of active subscribers.
func (b *Bus) Len() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.subs)
}
