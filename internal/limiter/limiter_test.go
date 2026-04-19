package limiter_test

import (
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/limiter"
	"github.com/user/portwatch/internal/scanner"
)

func makeEvent(t alert.EventType, port uint16) alert.Event {
	return alert.Event{
		Type: t,
		Port: scanner.Port{Number: port, Protocol: "tcp"},
	}
}

func TestNew_DefaultsToPositiveMax(t *testing.T) {
	l := limiter.New(0)
	if l.Count() != 0 {
		t.Fatalf("expected 0 count, got %d", l.Count())
	}
}

func TestApply_OpenedIncrementsCount(t *testing.T) {
	l := limiter.New(10)
	events := []alert.Event{makeEvent(alert.Opened, 80)}
	out, err := l.Apply(events)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("expected 1 event, got %d", len(out))
	}
	if l.Count() != 1 {
		t.Fatalf("expected count 1, got %d", l.Count())
	}
}

func TestApply_ClosedDecrementsCount(t *testing.T) {
	l := limiter.New(10)
	_, _ = l.Apply([]alert.Event{makeEvent(alert.Opened, 80)})
	out, err := l.Apply([]alert.Event{makeEvent(alert.Closed, 80)})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("expected 1 event, got %d", len(out))
	}
	if l.Count() != 0 {
		t.Fatalf("expected count 0, got %d", l.Count())
	}
}

func TestApply_ExceedsLimit_DropsEvent(t *testing.T) {
	l := limiter.New(1)
	_, _ = l.Apply([]alert.Event{makeEvent(alert.Opened, 80)})
	out, err := l.Apply([]alert.Event{makeEvent(alert.Opened, 443)})
	if err != limiter.ErrLimitExceeded {
		t.Fatalf("expected ErrLimitExceeded, got %v", err)
	}
	if len(out) != 0 {
		t.Fatalf("expected 0 allowed events, got %d", len(out))
	}
	if l.Count() != 1 {
		t.Fatalf("count should remain 1, got %d", l.Count())
	}
}

func TestReset_ClearsCount(t *testing.T) {
	l := limiter.New(10)
	_, _ = l.Apply([]alert.Event{makeEvent(alert.Opened, 80), makeEvent(alert.Opened, 443)})
	l.Reset()
	if l.Count() != 0 {
		t.Fatalf("expected 0 after reset, got %d", l.Count())
	}
}
