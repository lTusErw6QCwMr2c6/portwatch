package filter

import "github.com/user/portwatch/internal/scanner"

// Rule defines criteria for including or excluding ports from monitoring.
type Rule struct {
	PortStart uint16
	PortEnd   uint16
	Protocols []string // e.g. ["tcp", "udp"]; empty means all
	Exclude   bool     // if true, matching ports are excluded
}

// Filter applies a set of rules to a slice of ports.
type Filter struct {
	rules []Rule
}

// New creates a Filter with the given rules.
func New(rules []Rule) *Filter {
	return &Filter{rules: rules}
}

// Apply returns only the ports that pass all filter rules.
func (f *Filter) Apply(ports []scanner.Port) []scanner.Port {
	if len(f.rules) == 0 {
		return ports
	}
	out := make([]scanner.Port, 0, len(ports))
	for _, p := range ports {
		if f.allowed(p) {
			out = append(out, p)
		}
	}
	return out
}

// allowed returns true if the port is permitted by the filter rules.
func (f *Filter) allowed(p scanner.Port) bool {
	for _, r := range f.rules {
		if r.matchesPort(p) {
			return !r.Exclude
		}
	}
	// No rule matched — allow by default.
	return true
}

// matchesPort returns true when the rule applies to the given port.
func (r *Rule) matchesPort(p scanner.Port) bool {
	if p.Number < r.PortStart || p.Number > r.PortEnd {
		return false
	}
	if len(r.Protocols) == 0 {
		return true
	}
	for _, proto := range r.Protocols {
		if proto == p.Protocol {
			return true
		}
	}
	return false
}
