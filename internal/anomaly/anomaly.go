// Package anomaly detects unusual port activity based on rate and pattern heuristics.
package anomaly

import (
	"fmt"
	"sync"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// Severity represents the anomaly severity level.
type Severity string

const (
	SeverityLow    Severity = "low"
	SeverityMedium Severity = "medium"
	SeverityHigh   Severity = "high"
)

// Detection holds information about a detected anomaly.
type Detection struct {
	Event     alert.Event
	Reason    string
	Severity  Severity
	DetectedAt time.Time
}

func (d Detection) String() string {
	return fmt.Sprintf("[%s] %s — %s", d.Severity, d.Event, d.Reason)
}

// Detector evaluates events for anomalous behaviour.
type Detector struct {
	mu       sync.Mutex
	window   time.Duration
	threshold int
	buckets  map[string][]time.Time
}

// New returns a Detector with the given burst window and threshold.
func New(window time.Duration, threshold int) *Detector {
	return &Detector{
		window:    window,
		threshold: threshold,
		buckets:   make(map[string][]time.Time),
	}
}

// Evaluate checks whether the event constitutes an anomaly.
// Returns a Detection and true when an anomaly is found.
func (d *Detector) Evaluate(e alert.Event) (Detection, bool) {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := time.Now()
	key := fmt.Sprintf("%s:%d", e.Port.Protocol, e.Port.Number)

	times := d.buckets[key]
	cutoff := now.Add(-d.window)
	filtered := times[:0]
	for _, t := range times {
		if t.After(cutoff) {
			filtered = append(filtered, t)
		}
	}
	filtered = append(filtered, now)
	d.buckets[key] = filtered

	if len(filtered) >= d.threshold {
		sev := SeverityMedium
		if len(filtered) >= d.threshold*2 {
			sev = SeverityHigh
		}
		return Detection{
			Event:      e,
			Reason:     fmt.Sprintf("%d events on %s within %s", len(filtered), key, d.window),
			Severity:   sev,
			DetectedAt: now,
		}, true
	}
	return Detection{}, false
}

// Reset clears all tracked buckets.
func (d *Detector) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.buckets = make(map[string][]time.Time)
}
