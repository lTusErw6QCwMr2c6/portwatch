package cascade

import "github.com/user/portwatch/internal/alert"

// Chain composes multiple cascades into a single Stage.
// Events flow through each cascade in order.
func Chain(cascades ...*Cascade) Stage {
	return StageFunc(func(events []alert.Event) []alert.Event {
		current := events
		for _, c := range cascades {
			current = c.Run(current)
			if len(current) == 0 {
				return current
			}
		}
		return current
	})
}

// Merge combines results from multiple cascades running over the same
// input, returning a deduplicated union of their outputs.
func Merge(cascades ...*Cascade) Stage {
	return StageFunc(func(events []alert.Event) []alert.Event {
		seen := make(map[string]struct{})
		var out []alert.Event
		for _, c := range cascades {
			for _, e := range c.Run(events) {
				key := e.String()
				if _, ok := seen[key]; ok {
					continue
				}
				seen[key] = struct{}{}
				out = append(out, e)
			}
		}
		return out
	})
}
