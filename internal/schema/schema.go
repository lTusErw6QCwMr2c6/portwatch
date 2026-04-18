// Package schema validates event field structure before export or audit.
package schema

import (
	"errors"
	"fmt"

	"github.com/user/portwatch/internal/alert"
)

// Rule defines a validation constraint for an event field.
type Rule struct {
	Field    string
	Required bool
	AllowedValues []string
}

// Validator checks events against a set of rules.
type Validator struct {
	rules []Rule
}

// New returns a Validator with the given rules.
func New(rules []Rule) *Validator {
	return &Validator{rules: rules}
}

// DefaultRules returns a baseline set of validation rules for port events.
func DefaultRules() []Rule {
	return []Rule{
		{Field: "protocol", Required: true, AllowedValues: []string{"tcp", "udp"}},
		{Field: "action", Required: true, AllowedValues: []string{"opened", "closed"}},
	}
}

// Validate checks a single event against all rules.
// Returns a combined error if any rule is violated.
func (v *Validator) Validate(e alert.Event) error {
	var errs []error
	for _, r := range v.rules {
		if err := applyRule(r, e); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return errors.Join(errs...)
}

// ValidateAll validates a slice of events, returning all errors keyed by index.
func (v *Validator) ValidateAll(events []alert.Event) map[int]error {
	out := make(map[int]error)
	for i, e := range events {
		if err := v.Validate(e); err != nil {
			out[i] = err
		}
	}
	return out
}

func applyRule(r Rule, e alert.Event) error {
	var val string
	switch r.Field {
	case "protocol":
		val = e.Port.Protocol
	case "action":
		val = actionString(e)
	default:
		return fmt.Errorf("unknown field: %s", r.Field)
	}
	if r.Required && val == "" {
		return fmt.Errorf("field %q is required", r.Field)
	}
	if len(r.AllowedValues) > 0 && !contains(r.AllowedValues, val) {
		return fmt.Errorf("field %q value %q not in allowed set %v", r.Field, val, r.AllowedValues)
	}
	return nil
}

func actionString(e alert.Event) string {
	if e.Opened {
		return "opened"
	}
	return "closed"
}

func contains(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}
