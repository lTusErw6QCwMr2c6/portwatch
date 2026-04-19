package stale_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/stale"
)

func makeEvent(t alert.EventType, port int, proto string) alert.Event {
	return alert.Event{Type: t, Port: port, Protocol: proto}
}

func TestNew_NotNil(t *testing.T) {
	d := stale.New(time.Minute)
	if d == nil {
		t.Fatal("expected non-nil detector")
	}
}

func TestStale_NoEntries_ReturnsEmpty(t *testing.T) {
	d := stale.New(time.Millisecond)
	if got := d.Stale(); len(got) != 0 {
		t.Fatalf("expected empty, got %d", len(got))
	}
}

func TestTrack_OpenedThenStale(t *testing.T) {
	d := stale.New(10 * time.Millisecond)
	d.Track(makeEvent(alert.Opened, 8080, "tcp"))
	time.Sleep(20 * time.Millisecond)
	got := d.Stale()
	if len(got) != 1 {
		t.Fatalf("expected 1 stale entry, got %d", len(got))
	}
	if got[0].Port != 8080 {
		t.Errorf("expected port 8080, got %d", got[0].Port)
	}
}

func TestTrack_ClosedRemovesEntry(t *testing.T) {
	d := stale.New(10 * time.Millisecond)
	d.Track(makeEvent(alert.Opened, 9090, "tcp"))
	d.Track(makeEvent(alert.Closed, 9090, "tcp"))
	time.Sleep(20 * time.Millisecond)
	if got := d.Stale(); len(got) != 0 {
		t.Fatalf("expected 0 stale entries after close, got %d", len(got))
	}
}

func TestTrack_BelowMaxAge_NotStale(t *testing.T) {
	d := stale.New(time.Hour)
	d.Track(makeEvent(alert.Opened, 443, "tcp"))
	if got := d.Stale(); len(got) != 0 {
		t.Fatalf("expected 0 stale entries, got %d", len(got))
	}
}

func TestReset_ClearsEntries(t *testing.T) {
	d := stale.New(10 * time.Millisecond)
	d.Track(makeEvent(alert.Opened, 80, "tcp"))
	time.Sleep(20 * time.Millisecond)
	d.Reset()
	if got := d.Stale(); len(got) != 0 {
		t.Fatalf("expected 0 after reset, got %d", len(got))
	}
}
