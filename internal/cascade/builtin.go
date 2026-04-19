package cascade

import (
	"github.com/user/portwatch/internal/alert"
)

// DedupeStage removes duplicate events by port+protocol+type key.
func DedupeStage() Stage {
	return StageFunc(func(events []alert.Event) []alert.Event {
		seen := make(map[string]struct{}, len(events))
		out := events[:0:0]
		for _, e := range events {
			key := e.String()
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}
			out = append(out, e)
		}
		return out
	})
}

// LimitStage caps the number of events passed downstream.
func LimitStage(max int) Stage {
	return StageFunc(func(events []alert.Event) []alert.Event {
		if len(events) <= max {
			return events
		}
		return events[:max]
	})
}

// FilterTypeStage keeps only events matching the given type string.
func FilterTypeStage(typ string) Stage {
	return StageFunc(func(events []alert.Event) []alert.Event {
		out := events[:0:0]
		for _, e := range events {
			if e.Type == typ {
				out = append(out, e)
			}
		}
		return out
	})
}
