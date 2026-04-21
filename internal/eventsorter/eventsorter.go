package eventsorter

import (
	"sort"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// Field represents a sortable field on an event.
type Field int

const (
	ByTime Field = iota
	ByPort
	ByProtocol
)

// Order controls sort direction.
type Order int

const (
	Ascending Order = iota
	Descending
)

// Config holds sorter configuration.
type Config struct {
	Field Field
	Order Order
}

// Sorter sorts slices of alert.Event according to a configured field and order.
type Sorter struct {
	cfg Config
}

// New returns a Sorter with the given configuration.
func New(cfg Config) *Sorter {
	return &Sorter{cfg: cfg}
}

// Sort returns a new slice of events sorted according to the sorter's config.
// The original slice is not modified.
func (s *Sorter) Sort(events []alert.Event) []alert.Event {
	if len(events) == 0 {
		return []alert.Event{}
	}
	out := make([]alert.Event, len(events))
	copy(out, events)

	sort.SliceStable(out, func(i, j int) bool {
		less := s.less(out[i], out[j])
		if s.cfg.Order == Descending {
			return !less
		}
		return less
	})
	return out
}

func (s *Sorter) less(a, b alert.Event) bool {
	switch s.cfg.Field {
	case ByPort:
		return a.Port.Number < b.Port.Number
	case ByProtocol:
		return a.Port.Protocol < b.Port.Protocol
	case ByTime:
		fallthrough
	default:
		at, _ := time.Parse(time.RFC3339, a.Timestamp)
		bt, _ := time.Parse(time.RFC3339, b.Timestamp)
		return at.Before(bt)
	}
}
