package eventclassifier

import (
	"sync"

	"github.com/example/portwatch/internal/alert"
)

// Class represents a named classification label.
type Class string

const (
	ClassNormal    Class = "normal"
	ClassSuspicious Class = "suspicious"
	ClassCritical  Class = "critical"
	ClassUnknown   Class = "unknown"
)

// Rule maps a port range and protocol to a Class.
type Rule struct {
	PortStart uint16
	PortEnd   uint16
	Protocol  string // "tcp", "udp", or "" for any
	Class     Class
}

// Classifier assigns a Class to alert events based on registered rules.
type Classifier struct {
	mu    sync.RWMutex
	rules []Rule
}

// New returns a new Classifier with no rules.
func New() *Classifier {
	return &Classifier{}
}

// AddRule appends a classification rule.
func (c *Classifier) AddRule(r Rule) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.rules = append(c.rules, r)
}

// Classify returns the Class for the given event.
// Rules are evaluated in insertion order; the first match wins.
// Returns ClassUnknown if no rule matches.
func (c *Classifier) Classify(e alert.Event) Class {
	c.mu.RLock()
	defer c.mu.RUnlock()

	for _, r := range c.rules {
		if r.Protocol != "" && r.Protocol != e.Port.Protocol {
			continue
		}
		port := uint16(e.Port.Number)
		if port >= r.PortStart && port <= r.PortEnd {
			return r.Class
		}
	}
	return ClassUnknown
}

// ClassifyAll returns a map of event index to Class for a slice of events.
func (c *Classifier) ClassifyAll(events []alert.Event) map[int]Class {
	out := make(map[int]Class, len(events))
	for i, e := range events {
		out[i] = c.Classify(e)
	}
	return out
}

// Reset removes all registered rules.
func (c *Classifier) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.rules = c.rules[:0]
}
