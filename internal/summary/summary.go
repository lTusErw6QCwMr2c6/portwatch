package summary

import (
	"fmt"
	"strings"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/metrics"
)

// Report holds a periodic summary of port activity.
type Report struct {
	GeneratedAt time.Time
	Cycles      int
	Opened      int
	Closed      int
	Alerts      int
	TopEvents   []alert.Event
}

// Builder assembles summary reports from metrics and recent events.
type Builder struct {
	metrics *metrics.Metrics
	maxTop  int
}

// New returns a Builder using the provided metrics collector.
func New(m *metrics.Metrics, maxTop int) *Builder {
	if maxTop <= 0 {
		maxTop = 5
	}
	return &Builder{metrics: m, maxTop: maxTop}
}

// Build creates a Report from the current metrics state and recent events.
func (b *Builder) Build(events []alert.Event) Report {
	snap := b.metrics.Snapshot()
	top := events
	if len(top) > b.maxTop {
		top = top[:b.maxTop]
	}
	return Report{
		GeneratedAt: time.Now(),
		Cycles:      snap.Cycles,
		Opened:      snap.Opened,
		Closed:      snap.Closed,
		Alerts:      snap.Alerts,
		TopEvents:   top,
	}
}

// Format renders the report as a human-readable string.
func (r Report) Format() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("=== Port Activity Summary [%s] ===\n", r.GeneratedAt.Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("  Cycles : %d\n", r.Cycles))
	sb.WriteString(fmt.Sprintf("  Opened : %d\n", r.Opened))
	sb.WriteString(fmt.Sprintf("  Closed : %d\n", r.Closed))
	sb.WriteString(fmt.Sprintf("  Alerts : %d\n", r.Alerts))
	if len(r.TopEvents) > 0 {
		sb.WriteString("  Recent Events:\n")
		for _, e := range r.TopEvents {
			sb.WriteString(fmt.Sprintf("    - %s\n", e.String()))
		}
	}
	return sb.String()
}
