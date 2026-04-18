package trend

import (
	"testing"
	"time"
)

func TestRate_NoEvents_ReturnsZero(t *testing.T) {
	tr := New(time.Minute)
	if r := tr.Rate("tcp:8080"); r != 0 {
		t.Fatalf("expected 0, got %d", r)
	}
}

func TestRecord_IncreasesRate(t *testing.T) {
	tr := New(time.Minute)
	tr.Record("tcp:8080")
	tr.Record("tcp:8080")
	tr.Record("tcp:8080")
	if r := tr.Rate("tcp:8080"); r != 3 {
		t.Fatalf("expected 3, got %d", r)
	}
}

func TestRate_DifferentKeysAreIndependent(t *testing.T) {
	tr := New(time.Minute)
	tr.Record("tcp:80")
	tr.Record("udp:53")
	tr.Record("udp:53")
	if r := tr.Rate("tcp:80"); r != 1 {
		t.Fatalf("expected 1, got %d", r)
	}
	if r := tr.Rate("udp:53"); r != 2 {
		t.Fatalf("expected 2, got %d", r)
	}
}

func TestRate_EvictsExpiredEntries(t *testing.T) {
	tr := New(50 * time.Millisecond)
	tr.Record("tcp:443")
	time.Sleep(80 * time.Millisecond)
	if r := tr.Rate("tcp:443"); r != 0 {
		t.Fatalf("expected 0 after expiry, got %d", r)
	}
}

func TestReset_ClearsAll(t *testing.T) {
	tr := New(time.Minute)
	tr.Record("tcp:22")
	tr.Record("tcp:22")
	tr.Reset()
	if r := tr.Rate("tcp:22"); r != 0 {
		t.Fatalf("expected 0 after reset, got %d", r)
	}
}

func TestNew_DefaultWindow(t *testing.T) {
	tr := New(0)
	if tr.window != time.Minute {
		t.Fatalf("expected default window of 1m, got %v", tr.window)
	}
}
