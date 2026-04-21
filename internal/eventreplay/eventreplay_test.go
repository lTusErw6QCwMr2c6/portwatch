package eventreplay_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/eventreplay"
	"github.com/user/portwatch/internal/scanner"
)

func makeEvent(port int, t alert.EventType) alert.Event {
	return alert.Event{
		Type: t,
		Port: scanner.Port{Port: port, Protocol: "tcp"},
	}
}

func TestNew_DefaultCapacity(t *testing.T) {
	b := eventreplay.New(0)
	if b == nil {
		t.Fatal("expected non-nil buffer")
	}
	if b.Len() != 0 {
		t.Fatalf("expected empty buffer, got %d", b.Len())
	}
}

func TestAdd_StoresEntry(t *testing.T) {
	b := eventreplay.New(10)
	b.Add(makeEvent(8080, alert.EventOpened))
	if b.Len() != 1 {
		t.Fatalf("expected 1 entry, got %d", b.Len())
	}
}

func TestAdd_EvictsOldestWhenFull(t *testing.T) {
	b := eventreplay.New(3)
	for i := 0; i < 5; i++ {
		b.Add(makeEvent(1000+i, alert.EventOpened))
	}
	if b.Len() != 3 {
		t.Fatalf("expected 3 entries, got %d", b.Len())
	}
	all := b.All()
	// oldest two should have been evicted; first remaining port is 1002
	if all[0].Event.Port.Port != 1002 {
		t.Fatalf("expected port 1002, got %d", all[0].Event.Port.Port)
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	b := eventreplay.New(10)
	b.Add(makeEvent(443, alert.EventOpened))
	got := b.All()
	got[0].Event.Port.Port = 9999
	if b.All()[0].Event.Port.Port == 9999 {
		t.Fatal("All() should return a copy, not a reference")
	}
}

func TestSince_FiltersEntries(t *testing.T) {
	b := eventreplay.New(10)
	b.Add(makeEvent(80, alert.EventOpened))
	cutoff := time.Now()
	time.Sleep(2 * time.Millisecond)
	b.Add(makeEvent(443, alert.EventOpened))

	got := b.Since(cutoff)
	if len(got) != 1 {
		t.Fatalf("expected 1 entry since cutoff, got %d", len(got))
	}
	if got[0].Event.Port.Port != 443 {
		t.Fatalf("expected port 443, got %d", got[0].Event.Port.Port)
	}
}

func TestReset_ClearsBuffer(t *testing.T) {
	b := eventreplay.New(10)
	b.Add(makeEvent(22, alert.EventOpened))
	b.Reset()
	if b.Len() != 0 {
		t.Fatalf("expected empty buffer after reset, got %d", b.Len())
	}
}
