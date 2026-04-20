package eventchain

import (
	"fmt"
	"strings"

	"github.com/user/portwatch/internal/alert"
)

// Resolver walks a Chain and builds a human-readable ancestry path
// for a given event.
type Resolver struct {
	chain *Chain
}

// NewResolver returns a Resolver backed by c.
func NewResolver(c *Chain) *Resolver {
	return &Resolver{chain: c}
}

// Ancestry returns the chain of parent keys leading to event, from
// oldest ancestor to the event itself. If the event has no parent
// the slice contains only the event's own key.
func (r *Resolver) Ancestry(event alert.Event) []string {
	var path []string
	current := event.Port.String()
	seen := make(map[string]struct{})

	for {
		if _, dup := seen[current]; dup {
			break // cycle guard
		}
		seen[current] = struct{}{}
		path = append([]string{current}, path...)

		r.chain.mu.RLock()
		parent, ok := r.chain.parents[current]
		r.chain.mu.RUnlock()
		if !ok {
			break
		}
		current = parent
	}
	return path
}

// Format returns a single-line ancestry string, e.g. "80/tcp → 8080/tcp → 9090/tcp".
func (r *Resolver) Format(event alert.Event) string {
	parts := r.Ancestry(event)
	if len(parts) == 0 {
		return fmt.Sprintf("%s (no ancestry)", event.Port.String())
	}
	return strings.Join(parts, " → ")
}
