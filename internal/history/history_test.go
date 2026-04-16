package history

import (
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func makeEvent(port int, kind alert.EventKind) alert.Event {
	return alert.Event{
		Kind: kind,
		Port: scanner.Port{Number: port, Protocol: "tcp"},
	}
}

func TestNew_DefaultCapacity(t *testing.T) {
	h := New(0)
	if h.max != MaxEntries {
		t.Errorf("expected max %d, got %d", MaxEntries, h.max)
	}
}

func TestAdd_StoresEntry(t *testing.T) {
	h := New(10)
	h.Add(makeEvent(8080, alert.Opened))
	if h.Len() != 1 {
		t.Fatalf("expected 1 entry, got %d", h.Len())
	}
	entries := h.All()
	if entries[0].Event.Port.Number != 8080 {
		t.Errorf("unexpected port %d", entries[0].Event.Port.Number)
	}
}

func TestAdd_EvictsOldestWhenFull(t *testing.T) {
	h := New(3)
	for i := 1; i <= 4; i++ {
		h.Add(makeEvent(i, alert.Opened))
	}
	if h.Len() != 3 {
		t.Fatalf("expected 3 entries after eviction, got %d", h.Len())
	}
	entries := h.All()
	// oldest entry (port 1) should have been evicted
	if entries[0].Event.Port.Number != 2 {
		t.Errorf("expected oldest remaining port 2, got %d", entries[0].Event.Port.Number)
	}
	if entries[2].Event.Port.Number != 4 {
		t.Errorf("expected newest port 4, got %d", entries[2].Event.Port.Number)
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	h := New(10)
	h.Add(makeEvent(443, alert.Closed))
	a := h.All()
	a[0].Event.Port.Number = 9999
	b := h.All()
	if b[0].Event.Port.Number == 9999 {
		t.Error("All() should return a copy, not a reference")
	}
}

func TestClear_RemovesAllEntries(t *testing.T) {
	h := New(10)
	h.Add(makeEvent(80, alert.Opened))
	h.Add(makeEvent(443, alert.Opened))
	h.Clear()
	if h.Len() != 0 {
		t.Errorf("expected 0 entries after Clear, got %d", h.Len())
	}
}

func TestAdd_TimestampSet(t *testing.T) {
	h := New(5)
	h.Add(makeEvent(22, alert.Opened))
	entry := h.All()[0]
	if entry.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}
