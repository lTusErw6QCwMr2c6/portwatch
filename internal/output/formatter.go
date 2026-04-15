package output

import (
	"fmt"
	"strings"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// EventType represents the type of port change event.
type EventType string

const (
	EventOpened EventType = "OPENED"
	EventClosed EventType = "CLOSED"
)

// Event represents a single port change event.
type Event struct {
	Timestamp time.Time
	Type      EventType
	Port      scanner.FormatPort
	Protocol  string
	Process   string
}

// Format returns a human-readable string for the event.
func (e Event) Format() string {
	ts := e.Timestamp.Format("2006-01-02 15:04:05")
	process := e.Process
	if process == "" {
		process = "unknown"
	}
	return fmt.Sprintf("[%s] %-6s port %-6s proto=%-4s process=%s",
		ts, e.Type, string(e.Port), strings.ToLower(e.Protocol), process)
}

// FormatDiff converts scanner diff results into a slice of Events.
func FormatDiff(opened, closed []string, protocol string, ts time.Time) []Event {
	events := make([]Event, 0, len(opened)+len(closed))

	for _, p := range opened {
		events = append(events, Event{
			Timestamp: ts,
			Type:      EventOpened,
			Port:      scanner.FormatPort(p),
			Protocol:  protocol,
		})
	}

	for _, p := range closed {
		events = append(events, Event{
			Timestamp: ts,
			Type:      EventClosed,
			Port:      scanner.FormatPort(p),
			Protocol:  protocol,
		})
	}

	return events
}
