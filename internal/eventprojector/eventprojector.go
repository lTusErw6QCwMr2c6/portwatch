// Package eventprojector reduces a stream of alert events into a
// projected view keyed by port+protocol, retaining the most recent
// event for each unique endpoint.
package eventprojector

import (
	"sync"

	"github.com/user/portwatch/internal/alert"
)

// Projector maintains a latest-event projection keyed by port+protocol.
type Projector struct {
	mu      sync.RWMutex
	entries map[string]alert.Event
}

// New returns an initialised Projector.
func New() *Projector {
	return &Projector{
		entries: make(map[string]alert.Event),
	}
}

func key(e alert.Event) string {
	return e.Port.Protocol + ":" + e.Port.Address
}

// Apply upserts the event into the projection.
func (p *Projector) Apply(e alert.Event) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.entries[key(e)] = e
}

// Get returns the most recent event for the given protocol and address.
// ok is false when no entry exists.
func (p *Projector) Get(protocol, address string) (alert.Event, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	e, ok := p.entries[protocol+":"+address]
	return e, ok
}

// Snapshot returns a copy of all projected entries.
func (p *Projector) Snapshot() []alert.Event {
	p.mu.RLock()
	defer p.mu.RUnlock()
	out := make([]alert.Event, 0, len(p.entries))
	for _, e := range p.entries {
		out = append(out, e)
	}
	return out
}

// Remove deletes the projection entry for the given protocol and address.
func (p *Projector) Remove(protocol, address string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.entries, protocol+":"+address)
}

// Reset clears all projected entries.
func (p *Projector) Reset() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.entries = make(map[string]alert.Event)
}

// Len returns the number of projected entries.
func (p *Projector) Len() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return len(p.entries)
}
