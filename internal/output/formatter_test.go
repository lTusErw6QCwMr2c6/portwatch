package output

import (
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

var fixedTime = time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)

func TestEvent_Format_Opened(t *testing.T) {
	e := Event{
		Timestamp: fixedTime,
		Type:      EventOpened,
		Port:      scanner.FormatPort("8080"),
		Protocol:  "TCP",
		Process:   "nginx",
	}
	out := e.Format()
	if !strings.Contains(out, "OPENED") {
		t.Errorf("expected OPENED in output, got: %s", out)
	}
	if !strings.Contains(out, "8080") {
		t.Errorf("expected port 8080 in output, got: %s", out)
	}
	if !strings.Contains(out, "nginx") {
		t.Errorf("expected process nginx in output, got: %s", out)
	}
}

func TestEvent_Format_UnknownProcess(t *testing.T) {
	e := Event{
		Timestamp: fixedTime,
		Type:      EventClosed,
		Port:      scanner.FormatPort("443"),
		Protocol:  "TCP",
	}
	out := e.Format()
	if !strings.Contains(out, "unknown") {
		t.Errorf("expected 'unknown' process in output, got: %s", out)
	}
}

func TestFormatDiff_ReturnsCorrectCounts(t *testing.T) {
	opened := []string{"8080", "9090"}
	closed := []string{"3000"}

	events := FormatDiff(opened, closed, "TCP", fixedTime)

	if len(events) != 3 {
		t.Fatalf("expected 3 events, got %d", len(events))
	}

	openedCount := 0
	closedCount := 0
	for _, e := range events {
		switch e.Type {
		case EventOpened:
			openedCount++
		case EventClosed:
			closedCount++
		}
	}
	if openedCount != 2 {
		t.Errorf("expected 2 opened events, got %d", openedCount)
	}
	if closedCount != 1 {
		t.Errorf("expected 1 closed event, got %d", closedCount)
	}
}

func TestFormatDiff_Empty(t *testing.T) {
	events := FormatDiff(nil, nil, "UDP", fixedTime)
	if len(events) != 0 {
		t.Errorf("expected 0 events for empty diff, got %d", len(events))
	}
}
