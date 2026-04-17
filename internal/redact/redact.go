// Package redact provides utilities for scrubbing sensitive port/process
// information from log output before it is written or exported.
package redact

import (
	"regexp"
	"strings"
)

// Rule describes a single redaction pattern.
type Rule struct {
	Pattern     *regexp.Regexp
	Replacement string
}

// Redactor applies a set of Rules to strings.
type Redactor struct {
	rules []Rule
}

// New returns a Redactor with the supplied rules.
func New(rules []Rule) *Redactor {
	return &Redactor{rules: rules}
}

// DefaultRules returns a sensible baseline set of redaction rules that
// mask common sensitive tokens (passwords in connection strings, auth
// tokens, etc.) that may appear in process command lines.
func DefaultRules() []Rule {
	return []Rule{
		{
			Pattern:     regexp.MustCompile(`(?i)(password|passwd|secret|token)=[^\s&]+`),
			Replacement: "$1=***",
		},
		{
			Pattern:     regexp.MustCompile(`\b(?:\d{1,3}\.){3}\d{1,3}\b`),
			Replacement: "<ip>",
		},
	}
}

// Apply runs all rules against s and returns the sanitised result.
func (r *Redactor) Apply(s string) string {
	for _, rule := range r.rules {
		s = rule.Pattern.ReplaceAllString(s, rule.Replacement)
	}
	return s
}

// ApplyAll redacts every string in the slice, returning a new slice.
func (r *Redactor) ApplyAll(ss []string) []string {
	out := make([]string, len(ss))
	for i, s := range ss {
		out[i] = r.Apply(s)
	}
	return out
}

// ContainsSensitive returns true if s matches any registered pattern.
func (r *Redactor) ContainsSensitive(s string) bool {
	for _, rule := range r.rules {
		if rule.Pattern.MatchString(s) {
			return true
		}
	}
	return strings.ContainsAny(s, "\x00") // null bytes are always sensitive
}
