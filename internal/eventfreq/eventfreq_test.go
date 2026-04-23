package eventfreq_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/eventfreq"
	"github.com/user/portwatch/internal/scanner"
)

func makeEvent(port int, proto string) alert.Event {
	return alert.Event{
		Port: scanner.Port{Number: port, Protocol: proto},
		Type: alert.Opened,
	}
}

func TestNew_DefaultWindow(t *testing.T) {
	tr := eventfreq.New(0)
	if tr == nil {
		t.Fatal("expected non-nil tracker")
	}
}

func TestRecord_IncrementsCount(t *testing.T) {
	tr := eventfreq.New(5 * time.Second)
	e := makeEvent(8080, "tcp")

	if c := tr.Record(e); c != 1 {
		t.Fatalf("expected 1, got %d", c)
	}
	if c := tr.Record(e); c != 2 {
		t.Fatalf("expected 2, got %d", c)
	}
}

func TestCount_ReturnsCurrentCount(t *testing.T) {
	tr := eventfreq.New(5 * time.Second)
	e := makeEvent(443, "tcp")

	tr.Record(e)
	tr.Record(e)

	if c := tr.Count(e); c != 2 {
		t.Fatalf("expected 2, got %d", c)
	}
}

func TestCount_MissingKey_ReturnsZero(t *testing.T) {
	tr := eventfreq.New(5 * time.Second)
	e := makeEvent(9999, "udp")

	if c := tr.Count(e); c != 0 {
		t.Fatalf("expected 0, got %d", c)
	}
}

func TestRecord_DifferentKeysAreIndependent(t *testing.T) {
	tr := eventfreq.New(5 * time.Second)
	a := makeEvent(80, "tcp")
	b := makeEvent(80, "udp")

	tr.Record(a)
	tr.Record(a)
	tr.Record(b)

	if c := tr.Count(a); c != 2 {
		t.Fatalf("expected 2 for tcp, got %d", c)
	}
	if c := tr.Count(b); c != 1 {
		t.Fatalf("expected 1 for udp, got %d", c)
	}
}

func TestRecord_EvictsExpiredEntry(t *testing.T) {
	tr := eventfreq.New(10 * time.Millisecond)
	e := makeEvent(22, "tcp")

	tr.Record(e)
	time.Sleep(20 * time.Millisecond)

	// Recording again after expiry should reset count to 1.
	if c := tr.Record(e); c != 1 {
		t.Fatalf("expected count reset to 1, got %d", c)
	}
}

func TestReset_ClearsAll(t *testing.T) {
	tr := eventfreq.New(5 * time.Second)
	e := makeEvent(3306, "tcp")

	tr.Record(e)
	tr.Reset()

	if c := tr.Count(e); c != 0 {
		t.Fatalf("expected 0 after reset, got %d", c)
	}
}
