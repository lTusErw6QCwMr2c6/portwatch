package eventlabel

import (
	"errors"
	"sync"

	"github.com/user/portwatch/internal/alert"
)

// Labeler attaches static or dynamic string labels to events based on
// named rules. Each rule maps a label key to a value and an optional
// predicate that must match the event.
type Labeler struct {
	mu    sync.RWMutex
	rules map[string]rule
}

type rule struct {
	key       string
	value     string
	predicate func(alert.Event) bool
}

// New returns an empty Labeler.
func New() *Labeler {
	return &Labeler{rules: make(map[string]rule)}
}

// Register adds a labeling rule under name. predicate may be nil, in
// which case the label is applied unconditionally.
func (l *Labeler) Register(name, key, value string, predicate func(alert.Event) bool) error {
	if name == "" {
		return errors.New("eventlabel: rule name must not be empty")
	}
	if key == "" {
		return errors.New("eventlabel: label key must not be empty")
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if _, exists := l.rules[name]; exists {
		return errors.New("eventlabel: rule already registered: " + name)
	}
	l.rules[name] = rule{key: key, value: value, predicate: predicate}
	return nil
}

// Apply evaluates all registered rules against ev and returns a map of
// labels that matched. Returns nil when no rules match.
func (l *Labeler) Apply(ev alert.Event) map[string]string {
	l.mu.RLock()
	defer l.mu.RUnlock()
	var out map[string]string
	for _, r := range l.rules {
		if r.predicate != nil && !r.predicate(ev) {
			continue
		}
		if out == nil {
			out = make(map[string]string)
		}
		out[r.key] = r.value
	}
	return out
}

// Len returns the number of registered rules.
func (l *Labeler) Len() int {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return len(l.rules)
}
