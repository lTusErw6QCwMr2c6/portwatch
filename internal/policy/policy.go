package policy

import (
	"fmt"
	"strings"

	"github.com/user/portwatch/internal/alert"
)

// Action defines what to do when a rule matches.
type Action string

const (
	ActionAllow Action = "allow"
	ActionDeny  Action = "deny"
	ActionWarn  Action = "warn"
)

// Rule describes a policy rule applied to an alert event.
type Rule struct {
	Name     string
	Port     uint16
	Protocol string
	Action   Action
}

// Policy holds an ordered list of rules.
type Policy struct {
	rules []Rule
}

// New creates a Policy from the given rules.
func New(rules []Rule) *Policy {
	return &Policy{rules: rules}
}

// Evaluate returns the Action for the given event and the matching rule name.
// If no rule matches, ActionAllow is returned.
func (p *Policy) Evaluate(e alert.Event) (Action, string) {
	for _, r := range p.rules {
		if r.Port != 0 && r.Port != e.Port.Port {
			continue
		}
		if r.Protocol != "" && !strings.EqualFold(r.Protocol, e.Port.Protocol) {
			continue
		}
		return r.Action, r.Name
	}
	return ActionAllow, ""
}

// String returns a human-readable summary of the policy.
func (p *Policy) String() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Policy(%d rules)", len(p.rules))
	return sb.String()
}
