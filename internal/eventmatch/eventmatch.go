// Package eventmatch provides pattern-based matching against alert events.
package eventmatch

import (
	"strings"

	"github.com/user/portwatch/internal/alert"
)

// Rule defines a single match condition.
type Rule struct {
	Port     int
	Protocol string // "tcp", "udp", or "" for any
	Type     string // "opened", "closed", or "" for any
	Label    string // tag to attach on match
}

// Matcher evaluates events against a set of rules.
type Matcher struct {
	rules []Rule
}

// New returns a Matcher loaded with the given rules.
func New(rules []Rule) *Matcher {
	return &Matcher{rules: rules}
}

// Match returns the first Rule that matches the event, and true.
// If no rule matches, the zero Rule and false are returned.
func (m *Matcher) Match(e alert.Event) (Rule, bool) {
	for _, r := range m.rules {
		if r.Port != 0 && r.Port != e.Port.Port {
			continue
		}
		if r.Protocol != "" && !strings.EqualFold(r.Protocol, e.Port.Protocol) {
			continue
		}
		if r.Type != "" && !strings.EqualFold(r.Type, string(e.Type)) {
			continue
		}
		return r, true
	}
	return Rule{}, false
}

// MatchAll returns all rules that match the event.
func (m *Matcher) MatchAll(e alert.Event) []Rule {
	var matched []Rule
	for _, r := range m.rules {
		if r.Port != 0 && r.Port != e.Port.Port {
			continue
		}
		if r.Protocol != "" && !strings.EqualFold(r.Protocol, e.Port.Protocol) {
			continue
		}
		if r.Type != "" && !strings.EqualFold(r.Type, string(e.Type)) {
			continue
		}
		matched = append(matched, r)
	}
	return matched
}
