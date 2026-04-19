// Package tee fans out events to multiple sinks simultaneously.
package tee

import "github.com/user/portwatch/internal/alert"

// Sink is any type that can receive an alert event.
type Sink interface {
	Receive(event alert.Event) error
}

// Tee distributes a single event stream to multiple sinks.
type Tee struct {
	sinks []Sink
	stopOnErr bool
}

// New returns a Tee that writes to all provided sinks.
// If stopOnErr is true, the first error halts further delivery.
func New(stopOnErr bool, sinks ...Sink) *Tee {
	return &Tee{sinks: sinks, stopOnErr: stopOnErr}
}

// Add appends a sink to the tee.
func (t *Tee) Add(s Sink) {
	t.sinks = append(t.sinks, s)
}

// Len returns the number of registered sinks.
func (t *Tee) Len() int {
	return len(t.sinks)
}

// Send delivers event to every registered sink.
// Returns the first error encountered; remaining sinks are still called
// unless stopOnErr is true.
func (t *Tee) Send(event alert.Event) error {
	var first error
	for _, s := range t.sinks {
		if err := s.Receive(event); err != nil && first == nil {
			first = err
			if t.stopOnErr {
				return first
			}
		}
	}
	return first
}
