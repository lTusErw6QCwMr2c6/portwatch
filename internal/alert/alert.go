package alert

import (
	"fmt"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Level represents the severity of an alert.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelAlert Level = "ALERT"
)

// Event represents a single alert triggered by a port change.
type Event struct {
	Timestamp time.Time
	Level     Level
	Port      scanner.Port
	Opened    bool
}

// String returns a human-readable representation of the alert event.
func (e Event) String() string {
	action := "closed"
	if e.Opened {
		action = "opened"
	}
	return fmt.Sprintf("[%s] %s port %s %s",
		e.Timestamp.Format(time.RFC3339),
		e.Level,
		scanner.FormatPort(e.Port),
		action,
	)
}

// Threshold defines the criteria for escalating an alert level.
type Threshold struct {
	// WarnAfter escalates to WARN if more than this many ports change in one diff.
	WarnAfter int
	// AlertAfter escalates to ALERT if more than this many ports change in one diff.
	AlertAfter int
}

// DefaultThreshold returns a sensible default threshold configuration.
func DefaultThreshold() Threshold {
	return Threshold{
		WarnAfter:  5,
		AlertAfter: 20,
	}
}

// Evaluate determines the alert level based on the number of changed ports.
func Evaluate(changed int, t Threshold) Level {
	switch {
	case changed >= t.AlertAfter:
		return LevelAlert
	case changed >= t.WarnAfter:
		return LevelWarn
	default:
		return LevelInfo
	}
}

// BuildEvents creates a slice of alert Events from opened and closed port lists.
func BuildEvents(opened, closed []scanner.Port, t Threshold) []Event {
	total := len(opened) + len(closed)
	level := Evaluate(total, t)
	now := time.Now()

	events := make([]Event, 0, total)
	for _, p := range opened {
		events = append(events, Event{Timestamp: now, Level: level, Port: p, Opened: true})
	}
	for _, p := range closed {
		events = append(events, Event{Timestamp: now, Level: level, Port: p, Opened: false})
	}
	return events
}
