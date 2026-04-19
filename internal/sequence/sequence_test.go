package sequence_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/sequence"
)

func makeEvent(port int) alert.Event {
	return alert.Event{
		Port:   scanner.Port{Number: port, Protocol: "tcp"},
		Action: alert.Opened,
	}
}

func TestNew_NotNil(t *testing.T) {
	tr := sequence.New(time.Second, 10)
	if tr == nil {
		t.Fatal("expected non-nil tracker")
	}
}

func TestAdd_And_Count(t *testing.T) {
	tr := sequence.New(time.Second, 10)
	tr.Add(makeEvent(80))
	tr.Add(makeEvent(443))
	if got := tr.Count(); got != 2 {
		t.Fatalf("expected 2, got %d", got)
	}
}

func TestRecent_ReturnsEntries(t *testing.T) {
	tr := sequence.New(time.Second, 10)
	tr.Add(makeEvent(22))
	entries := tr.Recent()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Event.Port.Number != 22 {
		t.Errorf("unexpected port %d", entries[0].Event.Port.Number)
	}
}

func TestEviction_RemovesExpiredEntries(t *testing.T) {
	tr := sequence.New(50*time.Millisecond, 10)
	tr.Add(makeEvent(8080))
	time.Sleep(80 * time.Millisecond)
	if got := tr.Count(); got != 0 {
		t.Fatalf("expected 0 after expiry, got %d", got)
	}
}

func TestReset_ClearsEntries(t *testing.T) {
	tr := sequence.New(time.Second, 10)
	tr.Add(makeEvent(3000))
	tr.Reset()
	if got := tr.Count(); got != 0 {
		t.Fatalf("expected 0 after reset, got %d", got)
	}
}

func TestMaxSize_EvictsOldest(t *testing.T) {
	tr := sequence.New(time.Minute, 3)
	for i := 0; i < 5; i++ {
		tr.Add(makeEvent(i))
	}
	if got := tr.Count(); got != 3 {
		t.Fatalf("expected 3 (maxSize), got %d", got)
	}
}
