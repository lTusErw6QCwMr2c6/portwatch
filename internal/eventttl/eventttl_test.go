package eventttl_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/eventttl"
	"github.com/user/portwatch/internal/scanner"
)

func makeEvent(port int) alert.Event {
	return alert.Event{
		Type: alert.Opened,
		Port: scanner.Port{Number: port, Protocol: "tcp"},
	}
}

func TestAdd_And_Get_LiveEntry(t *testing.T) {
	s := eventttl.New(5 * time.Second)
	ev := makeEvent(8080)
	s.Add("k1", ev)

	entry, ok := s.Get("k1")
	if !ok {
		t.Fatal("expected entry to be present")
	}
	if entry.Event.Port.Number != 8080 {
		t.Errorf("unexpected port: %d", entry.Event.Port.Number)
	}
}

func TestGet_MissingKey_ReturnsFalse(t *testing.T) {
	s := eventttl.New(5 * time.Second)
	_, ok := s.Get("missing")
	if ok {
		t.Error("expected false for missing key")
	}
}

func TestGet_ExpiredEntry_ReturnsFalse(t *testing.T) {
	s := eventttl.New(1 * time.Millisecond)
	s.Add("k1", makeEvent(9090))
	time.Sleep(5 * time.Millisecond)

	_, ok := s.Get("k1")
	if ok {
		t.Error("expected expired entry to return false")
	}
}

func TestEvict_RemovesExpiredEntries(t *testing.T) {
	s := eventttl.New(1 * time.Millisecond)
	s.Add("k1", makeEvent(80))
	s.Add("k2", makeEvent(443))
	time.Sleep(5 * time.Millisecond)

	evicted := s.Evict()
	if len(evicted) != 2 {
		t.Errorf("expected 2 evicted, got %d", len(evicted))
	}
	if s.Len() != 0 {
		t.Errorf("expected store to be empty after eviction, got %d", s.Len())
	}
}

func TestEvict_PreservesLiveEntries(t *testing.T) {
	s := eventttl.New(5 * time.Second)
	s.Add("live", makeEvent(22))

	evicted := s.Evict()
	if len(evicted) != 0 {
		t.Errorf("expected 0 evicted, got %d", len(evicted))
	}
	if s.Len() != 1 {
		t.Errorf("expected 1 live entry, got %d", s.Len())
	}
}

func TestReset_ClearsAllEntries(t *testing.T) {
	s := eventttl.New(5 * time.Second)
	s.Add("a", makeEvent(1))
	s.Add("b", makeEvent(2))
	s.Reset()

	if s.Len() != 0 {
		t.Errorf("expected 0 after reset, got %d", s.Len())
	}
}

func TestAdd_RefreshesExpiry(t *testing.T) {
	s := eventttl.New(50 * time.Millisecond)
	s.Add("k", makeEvent(3306))
	time.Sleep(30 * time.Millisecond)
	s.Add("k", makeEvent(3306)) // refresh
	time.Sleep(30 * time.Millisecond)

	_, ok := s.Get("k")
	if !ok {
		t.Error("expected refreshed entry to still be live")
	}
}
