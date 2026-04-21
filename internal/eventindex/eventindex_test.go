package eventindex

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func makeEvent(port int, proto, kind string) alert.Event {
	return alert.Event{
		Kind: alert.Kind(kind),
		Port: scanner.Port{
			Number:   port,
			Protocol: proto,
		},
	}
}

func TestNew_NotNil(t *testing.T) {
	idx := New(time.Minute)
	if idx == nil {
		t.Fatal("expected non-nil index")
	}
}

func TestAdd_And_ByPort(t *testing.T) {
	idx := New(time.Minute)
	ev := makeEvent(8080, "tcp", "opened")
	idx.Add(ev)

	results := idx.ByPort(8080)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Port.Number != 8080 {
		t.Errorf("expected port 8080, got %d", results[0].Port.Number)
	}
}

func TestByPort_MissingPort_ReturnsEmpty(t *testing.T) {
	idx := New(time.Minute)
	results := idx.ByPort(9999)
	if len(results) != 0 {
		t.Errorf("expected empty, got %d results", len(results))
	}
}

func TestByProtocol_ReturnsMatchingEvents(t *testing.T) {
	idx := New(time.Minute)
	idx.Add(makeEvent(80, "tcp", "opened"))
	idx.Add(makeEvent(53, "udp", "opened"))
	idx.Add(makeEvent(443, "tcp", "closed"))

	results := idx.ByProtocol("tcp")
	if len(results) != 2 {
		t.Errorf("expected 2 tcp events, got %d", len(results))
	}
}

func TestByType_ReturnsMatchingEvents(t *testing.T) {
	idx := New(time.Minute)
	idx.Add(makeEvent(80, "tcp", "opened"))
	idx.Add(makeEvent(443, "tcp", "closed"))
	idx.Add(makeEvent(8080, "tcp", "opened"))

	results := idx.ByType("opened")
	if len(results) != 2 {
		t.Errorf("expected 2 opened events, got %d", len(results))
	}
}

func TestEvict_RemovesExpiredEntries(t *testing.T) {
	idx := New(10 * time.Millisecond)
	idx.Add(makeEvent(8080, "tcp", "opened"))

	time.Sleep(30 * time.Millisecond)
	idx.Evict()

	results := idx.ByPort(8080)
	if len(results) != 0 {
		t.Errorf("expected 0 results after eviction, got %d", len(results))
	}
}

func TestNew_DefaultTTL_WhenZero(t *testing.T) {
	idx := New(0)
	if idx.ttl != 5*time.Minute {
		t.Errorf("expected default TTL of 5m, got %v", idx.ttl)
	}
}
