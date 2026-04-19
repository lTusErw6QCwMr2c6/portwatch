package cascade_test

import (
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/cascade"
	"github.com/user/portwatch/internal/scanner"
)

func makeEvents(types ...string) []alert.Event {
	events := make([]alert.Event, len(types))
	for i, t := range types {
		events[i] = alert.Event{
			Type: t,
			Port: scanner.Port{Port: uint16(8000 + i), Protocol: "tcp"},
		}
	}
	return events
}

func TestNew_NotNil(t *testing.T) {
	c := cascade.New()
	if c == nil {
		t.Fatal("expected non-nil cascade")
	}
}

func TestRun_NoStages_ReturnsAll(t *testing.T) {
	c := cascade.New()
	events := makeEvents("opened", "closed")
	out := c.Run(events)
	if len(out) != 2 {
		t.Fatalf("expected 2 events, got %d", len(out))
	}
}

func TestRun_EmptyInput_ReturnsEmpty(t *testing.T) {
	c := cascade.New(cascade.LimitStage(10))
	out := c.Run(nil)
	if len(out) != 0 {
		t.Fatalf("expected 0 events, got %d", len(out))
	}
}

func TestLimitStage_CapsOutput(t *testing.T) {
	c := cascade.New(cascade.LimitStage(1))
	events := makeEvents("opened", "opened", "closed")
	out := c.Run(events)
	if len(out) != 1 {
		t.Fatalf("expected 1 event, got %d", len(out))
	}
}

func TestDedupeStage_RemovesDuplicates(t *testing.T) {
	events := []alert.Event{
		{Type: "opened", Port: scanner.Port{Port: 80, Protocol: "tcp"}},
		{Type: "opened", Port: scanner.Port{Port: 80, Protocol: "tcp"}},
	}
	c := cascade.New(cascade.DedupeStage())
	out := c.Run(events)
	if len(out) != 1 {
		t.Fatalf("expected 1 unique event, got %d", len(out))
	}
}

func TestFilterTypeStage_KeepsMatchingType(t *testing.T) {
	c := cascade.New(cascade.FilterTypeStage("opened"))
	events := makeEvents("opened", "closed", "opened")
	out := c.Run(events)
	if len(out) != 2 {
		t.Fatalf("expected 2 opened events, got %d", len(out))
	}
}

func TestAdd_IncreasesLen(t *testing.T) {
	c := cascade.New()
	c.Add(cascade.LimitStage(5))
	if c.Len() != 1 {
		t.Fatalf("expected len 1, got %d", c.Len())
	}
}
