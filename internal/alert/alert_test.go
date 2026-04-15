package alert

import (
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func TestDefaultThreshold(t *testing.T) {
	th := DefaultThreshold()
	if th.WarnAfter != 5 {
		t.Errorf("expected WarnAfter=5, got %d", th.WarnAfter)
	}
	if th.AlertAfter != 20 {
		t.Errorf("expected AlertAfter=20, got %d", th.AlertAfter)
	}
}

func TestEvaluate(t *testing.T) {
	th := DefaultThreshold()
	tests := []struct {
		changed  int
		expected Level
	}{
		{0, LevelInfo},
		{4, LevelInfo},
		{5, LevelWarn},
		{19, LevelWarn},
		{20, LevelAlert},
		{100, LevelAlert},
	}
	for _, tc := range tests {
		got := Evaluate(tc.changed, th)
		if got != tc.expected {
			t.Errorf("Evaluate(%d): expected %s, got %s", tc.changed, tc.expected, got)
		}
	}
}

func TestEvent_String_Opened(t *testing.T) {
	e := Event{
		Timestamp: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		Level:     LevelInfo,
		Port:      scanner.Port{Port: 8080, Protocol: "tcp"},
		Opened:    true,
	}
	result := e.String()
	if !strings.Contains(result, "opened") {
		t.Errorf("expected 'opened' in output, got: %s", result)
	}
	if !strings.Contains(result, "INFO") {
		t.Errorf("expected 'INFO' in output, got: %s", result)
	}
}

func TestEvent_String_Closed(t *testing.T) {
	e := Event{
		Timestamp: time.Now(),
		Level:     LevelWarn,
		Port:      scanner.Port{Port: 443, Protocol: "tcp"},
		Opened:    false,
	}
	result := e.String()
	if !strings.Contains(result, "closed") {
		t.Errorf("expected 'closed' in output, got: %s", result)
	}
	if !strings.Contains(result, "WARN") {
		t.Errorf("expected 'WARN' in output, got: %s", result)
	}
}

func TestBuildEvents(t *testing.T) {
	opened := []scanner.Port{{Port: 8080, Protocol: "tcp"}, {Port: 9090, Protocol: "udp"}}
	closed := []scanner.Port{{Port: 22, Protocol: "tcp"}}
	th := DefaultThreshold()

	events := BuildEvents(opened, closed, th)
	if len(events) != 3 {
		t.Fatalf("expected 3 events, got %d", len(events))
	}
	if !events[0].Opened {
		t.Error("expected first event to be opened")
	}
	if events[2].Opened {
		t.Error("expected last event to be closed")
	}
}
