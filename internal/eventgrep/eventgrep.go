package eventgrep

import (
	"regexp"
	"strings"

	"github.com/user/portwatch/internal/alert"
)

// Matcher defines a compiled grep rule used to match against event fields.
type Matcher struct {
	pattern *regexp.Regexp
	fields  []string
}

// Grepper searches events using pattern-based rules.
type Grepper struct {
	matchers []*Matcher
}

// New creates a Grepper with no rules.
func New() *Grepper {
	return &Grepper{}
}

// Add compiles and registers a pattern to match against the given event fields.
// If fields is empty, all supported fields are searched.
func (g *Grepper) Add(pattern string, fields ...string) error {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return err
	}
	if len(fields) == 0 {
		fields = []string{"type", "protocol", "process"}
	}
	g.matchers = append(g.matchers, &Matcher{pattern: re, fields: fields})
	return nil
}

// Match returns true if the event matches any registered pattern.
func (g *Grepper) Match(e alert.Event) bool {
	for _, m := range g.matchers {
		if matchEvent(m, e) {
			return true
		}
	}
	return false
}

// Filter returns only the events that match at least one registered pattern.
func (g *Grepper) Filter(events []alert.Event) []alert.Event {
	out := make([]alert.Event, 0, len(events))
	for _, e := range events {
		if g.Match(e) {
			out = append(out, e)
		}
	}
	return out
}

// Len returns the number of registered matchers.
func (g *Grepper) Len() int { return len(g.matchers) }

func matchEvent(m *Matcher, e alert.Event) bool {
	for _, field := range m.fields {
		var val string
		switch strings.ToLower(field) {
		case "type":
			val = string(e.Type)
		case "protocol":
			val = e.Port.Protocol
		case "process":
			val = e.Port.Process
		}
		if m.pattern.MatchString(val) {
			return true
		}
	}
	return false
}
