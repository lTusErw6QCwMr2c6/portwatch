package eventsorter_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/eventsorter"
	"github.com/user/portwatch/internal/scanner"
)

func makeEvent(port int, proto string, ts time.Time) alert.Event {
	return alert.Event{
		Port:      scanner.Port{Number: port, Protocol: proto},
		Timestamp: ts.Format(time.RFC3339),
	}
}

func TestSort_ByPort_Ascending(t *testing.T) {
	now := time.Now()
	events := []alert.Event{
		makeEvent(8080, "tcp", now),
		makeEvent(22, "tcp", now),
		makeEvent(443, "tcp", now),
	}
	s := eventsorter.New(eventsorter.Config{Field: eventsorter.ByPort, Order: eventsorter.Ascending})
	out := s.Sort(events)
	if out[0].Port.Number != 22 || out[1].Port.Number != 443 || out[2].Port.Number != 8080 {
		t.Errorf("unexpected order: %v", out)
	}
}

func TestSort_ByPort_Descending(t *testing.T) {
	now := time.Now()
	events := []alert.Event{
		makeEvent(22, "tcp", now),
		makeEvent(8080, "tcp", now),
		makeEvent(443, "tcp", now),
	}
	s := eventsorter.New(eventsorter.Config{Field: eventsorter.ByPort, Order: eventsorter.Descending})
	out := s.Sort(events)
	if out[0].Port.Number != 8080 || out[1].Port.Number != 443 || out[2].Port.Number != 22 {
		t.Errorf("unexpected order: %v", out)
	}
}

func TestSort_ByTime_Ascending(t *testing.T) {
	base := time.Now()
	events := []alert.Event{
		makeEvent(80, "tcp", base.Add(2*time.Second)),
		makeEvent(443, "tcp", base),
		makeEvent(22, "tcp", base.Add(time.Second)),
	}
	s := eventsorter.New(eventsorter.Config{Field: eventsorter.ByTime, Order: eventsorter.Ascending})
	out := s.Sort(events)
	if out[0].Port.Number != 443 || out[1].Port.Number != 22 || out[2].Port.Number != 80 {
		t.Errorf("unexpected time order: %v", out)
	}
}

func TestSort_ByProtocol_Ascending(t *testing.T) {
	now := time.Now()
	events := []alert.Event{
		makeEvent(80, "udp", now),
		makeEvent(443, "tcp", now),
	}
	s := eventsorter.New(eventsorter.Config{Field: eventsorter.ByProtocol, Order: eventsorter.Ascending})
	out := s.Sort(events)
	if out[0].Port.Protocol != "tcp" {
		t.Errorf("expected tcp first, got %s", out[0].Port.Protocol)
	}
}

func TestSort_Empty_ReturnsEmpty(t *testing.T) {
	s := eventsorter.New(eventsorter.Config{Field: eventsorter.ByPort})
	out := s.Sort([]alert.Event{})
	if len(out) != 0 {
		t.Errorf("expected empty slice, got %d elements", len(out))
	}
}

func TestSort_DoesNotMutateOriginal(t *testing.T) {
	now := time.Now()
	events := []alert.Event{
		makeEvent(8080, "tcp", now),
		makeEvent(22, "tcp", now),
	}
	orig := events[0].Port.Number
	s := eventsorter.New(eventsorter.Config{Field: eventsorter.ByPort, Order: eventsorter.Ascending})
	s.Sort(events)
	if events[0].Port.Number != orig {
		t.Error("original slice was mutated")
	}
}
