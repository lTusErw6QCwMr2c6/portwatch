// Package fence provides a port boundary guard that restricts
// monitoring to a declared set of allowed port ranges.
package fence

import "fmt"

// Range represents an inclusive port range with an optional protocol filter.
type Range struct {
	Start    int
	End      int
	Protocol string // "tcp", "udp", or "" for both
}

// Fence holds a set of allowed ranges and decides whether a port is in scope.
type Fence struct {
	ranges []Range
}

// New creates a Fence from the provided ranges.
func New(ranges []Range) (*Fence, error) {
	for _, r := range ranges {
		if r.Start < 1 || r.End > 65535 || r.Start > r.End {
			return nil, fmt.Errorf("fence: invalid range %d-%d", r.Start, r.End)
		}
	}
	return &Fence{ranges: ranges}, nil
}

// Allow returns true if the given port and protocol fall within any declared range.
func (f *Fence) Allow(port int, protocol string) bool {
	for _, r := range f.ranges {
		if port < r.Start || port > r.End {
			continue
		}
		if r.Protocol != "" && r.Protocol != protocol {
			continue
		}
		return true
	}
	return false
}

// Ranges returns a copy of the configured ranges.
func (f *Fence) Ranges() []Range {
	out := make([]Range, len(f.ranges))
	copy(out, f.ranges)
	return out
}

// Len returns the number of configured ranges.
func (f *Fence) Len() int { return len(f.ranges) }
