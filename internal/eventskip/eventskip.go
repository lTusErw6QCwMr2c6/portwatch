// Package eventskip provides a conditional skip mechanism that allows
// events to be selectively bypassed based on registered predicate rules.
package eventskip

import (
	"errors"
	"sync"

	"github.com/user/portwatch/internal/alert"
)

// Predicate is a function that returns true if an event should be skipped.
type Predicate func(e alert.Event) bool

type rule struct {
	name      string
	predicate Predicate
}

// Skipper holds named skip rules evaluated in registration order.
type Skipper struct {
	mu    sync.RWMutex
	rules []rule
}

// New returns an empty Skipper.
func New() *Skipper {
	return &Skipper{}
}

// Register adds a named predicate rule. Returns an error if the name is
// empty or already registered.
func (s *Skipper) Register(name string, p Predicate) error {
	if name == "" {
		return errors.New("eventskip: rule name must not be empty")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, r := range s.rules {
		if r.name == name {
			return errors.New("eventskip: rule already registered: " + name)
		}
	}
	s.rules = append(s.rules, rule{name: name, predicate: p})
	return nil
}

// Deregister removes a rule by name. Returns false if not found.
func (s *Skipper) Deregister(name string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, r := range s.rules {
		if r.name == name {
			s.rules = append(s.rules[:i], s.rules[i+1:]...)
			return true
		}
	}
	return false
}

// Skip returns true if any registered predicate matches the event.
func (s *Skipper) Skip(e alert.Event) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, r := range s.rules {
		if r.predicate(e) {
			return true
		}
	}
	return false
}

// Len returns the number of registered rules.
func (s *Skipper) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.rules)
}
