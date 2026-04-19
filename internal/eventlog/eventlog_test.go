package eventlog

import (
	"fmt"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func makeEntry(id string, port int, t time.Time) Entry {
	return Entry{
		ID:        id,
		Timestamp: t,
		Event: alert.Event{
			Type: alert.Opened,
			Port: scanner.Port{Number: port, Protocol: "tcp"},
		},
		Source: "test",
	}
}

func TestNew_DefaultCapacity(t *testing.T) {
	l := New(0)
	if l.capacity != 256 {
		t.Fatalf("expected 256, got %d", l.capacity)
	}
}

func TestAppend_StoresEntry(t *testing.T) {
	l := New(10)
	e := makeEntry("e1", 8080, time.Now())
	l.Append(e)
	if l.Len() != 1 {
		t.Fatalf("expected 1 entry, got %d", l.Len())
	}
}

func TestAppend_EvictsOldestWhenFull(t *testing.T) {
	l := New(3)
	for i := 0; i < 4; i++ {
		l.Append(makeEntry(fmt.Sprintf("e%d", i), i+1, time.Now()))
	}
	if l.Len() != 3 {
		t.Fatalf("expected 3 entries, got %d", l.Len())
	}
	if l.All()[0].ID != "e1" {
		t.Fatalf("expected oldest evicted, got %s", l.All()[0].ID)
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	l := New(10)
	l.Append(makeEntry("e1", 80, time.Now()))
	all := l.All()
	all[0].ID = "mutated"
	if l.All()[0].ID == "mutated" {
		t.Fatal("All() should return a copy")
	}
}

func TestSince_FiltersOldEntries(t *testing.T) {
	l := New(10)
	now := time.Now()
	l.Append(makeEntry("old", 80, now.Add(-2*time.Minute)))
	l.Append(makeEntry("new", 443, now.Add(time.Second)))
	result := l.Since(now)
	if len(result) != 1 || result[0].ID != "new" {
		t.Fatalf("expected 1 new entry, got %+v", result)
	}
}

func TestClear_RemovesAll(t *testing.T) {
	l := New(10)
	l.Append(makeEntry("e1", 80, time.Now()))
	l.Clear()
	if l.Len() != 0 {
		t.Fatal("expected empty log after Clear")
	}
}

func TestEntry_String_ContainsID(t *testing.T) {
	e := makeEntry("abc123", 80, time.Now())
	if s := e.String(); len(s) == 0 {
		t.Fatal("String() returned empty")
	}
}
