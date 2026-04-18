package buffer_test

import (
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/buffer"
	"github.com/user/portwatch/internal/scanner"
)

func makeEvent(port uint16) alert.Event {
	return alert.Event{
		Port: scanner.Port{Number: port, Protocol: "tcp"},
		Type: alert.Opened,
	}
}

func TestNew_DefaultCapacity(t *testing.T) {
	b := buffer.New(0)
	if b == nil {
		t.Fatal("expected non-nil buffer")
	}
}

func TestAdd_StoresEvent(t *testing.T) {
	b := buffer.New(4)
	b.Add(makeEvent(80))
	if b.Len() != 1 {
		t.Fatalf("expected 1, got %d", b.Len())
	}
}

func TestAll_ReturnsInsertionOrder(t *testing.T) {
	b := buffer.New(4)
	ports := []uint16{80, 443, 8080}
	for _, p := range ports {
		b.Add(makeEvent(p))
	}
	events := b.All()
	for i, e := range events {
		if e.Port.Number != ports[i] {
			t.Errorf("index %d: expected port %d, got %d", i, ports[i], e.Port.Number)
		}
	}
}

func TestAdd_OverwritesOldestWhenFull(t *testing.T) {
	b := buffer.New(3)
	for _, p := range []uint16{1, 2, 3, 4} {
		b.Add(makeEvent(p))
	}
	events := b.All()
	if len(events) != 3 {
		t.Fatalf("expected 3 events, got %d", len(events))
	}
	if events[0].Port.Number != 2 {
		t.Errorf("expected oldest to be 2, got %d", events[0].Port.Number)
	}
	if events[2].Port.Number != 4 {
		t.Errorf("expected newest to be 4, got %d", events[2].Port.Number)
	}
}

func TestReset_ClearsBuffer(t *testing.T) {
	b := buffer.New(4)
	b.Add(makeEvent(22))
	b.Add(makeEvent(80))
	b.Reset()
	if b.Len() != 0 {
		t.Fatalf("expected 0 after reset, got %d", b.Len())
	}
}
