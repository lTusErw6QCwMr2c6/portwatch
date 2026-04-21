package eventscope

import (
	"fmt"
	"sync"

	"github.com/user/portwatch/internal/alert"
)

// Scope defines a named boundary that restricts event processing to a
// specific port range and protocol combination.
type Scope struct {
	Name     string
	Protocol string
	PortMin  int
	PortMax  int
}

// Scoper holds a collection of named scopes and matches events against them.
type Scoper struct {
	mu     sync.RWMutex
	scopes map[string]Scope
}

// New returns an empty Scoper.
func New() *Scoper {
	return &Scoper{
		scopes: make(map[string]Scope),
	}
}

// Add registers a scope by name. Returns an error if the port range is invalid
// or a scope with that name already exists.
func (s *Scoper) Add(sc Scope) error {
	if sc.PortMin < 0 || sc.PortMax < sc.PortMin {
		return fmt.Errorf("eventscope: invalid port range %d-%d", sc.PortMin, sc.PortMax)
	}
	if sc.Name == "" {
		return fmt.Errorf("eventscope: scope name must not be empty")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.scopes[sc.Name]; exists {
		return fmt.Errorf("eventscope: scope %q already registered", sc.Name)
	}
	s.scopes[sc.Name] = sc
	return nil
}

// Remove deletes a scope by name.
func (s *Scoper) Remove(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.scopes, name)
}

// Match returns the names of all scopes that match the given event.
// A scope matches when the event port falls within [PortMin, PortMax] and
// the protocol matches (or the scope protocol is empty, meaning any).
func (s *Scoper) Match(ev alert.Event) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var matched []string
	for _, sc := range s.scopes {
		if ev.Port < sc.PortMin || ev.Port > sc.PortMax {
			continue
		}
		if sc.Protocol != "" && sc.Protocol != ev.Protocol {
			continue
		}
		matched = append(matched, sc.Name)
	}
	return matched
}

// Len returns the number of registered scopes.
func (s *Scoper) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.scopes)
}
