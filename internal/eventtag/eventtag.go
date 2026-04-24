package eventtag

import (
	"fmt"
	"sync"

	"github.com/user/portwatch/internal/alert"
)

// Tagger attaches string labels to events based on matching rules.
type Tagger struct {
	mu    sync.RWMutex
	rules []rule
}

type rule struct {
	name     string
	port     int
	protocol string // empty matches any
	tags     []string
}

// New returns an empty Tagger.
func New() *Tagger {
	return &Tagger{}
}

// Register adds a tagging rule. name must be unique and non-empty.
func (t *Tagger) Register(name string, port int, protocol string, tags []string) error {
	if name == "" {
		return fmt.Errorf("eventtag: rule name must not be empty")
	}
	if len(tags) == 0 {
		return fmt.Errorf("eventtag: rule %q must specify at least one tag", name)
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	for _, r := range t.rules {
		if r.name == name {
			return fmt.Errorf("eventtag: rule %q already registered", name)
		}
	}
	t.rules = append(t.rules, rule{name: name, port: port, protocol: protocol, tags: tags})
	return nil
}

// Apply returns all tags that match the given event. Returns nil if no rules match.
func (t *Tagger) Apply(ev alert.Event) []string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	var out []string
	for _, r := range t.rules {
		if r.port != ev.Port {
			continue
		}
		if r.protocol != "" && r.protocol != ev.Protocol {
			continue
		}
		out = append(out, r.tags...)
	}
	return out
}

// Remove deletes a rule by name. Returns false if not found.
func (t *Tagger) Remove(name string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	for i, r := range t.rules {
		if r.name == name {
			t.rules = append(t.rules[:i], t.rules[i+1:]...)
			return true
		}
	}
	return false
}

// Len returns the number of registered rules.
func (t *Tagger) Len() int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return len(t.rules)
}
