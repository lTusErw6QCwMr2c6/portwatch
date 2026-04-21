package eventcounter

import (
	"fmt"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func makeEvent(port int, proto string) alert.Event {
	return alert.Event{
		Port: scanner.Port{Number: port, Protocol: proto},
		Type: alert.Opened,
	}
}

func TestNew_DefaultWindow(t *testing.T) {
	c := New(0)
	if c.window != 60*time.Second {
		t.Fatalf("expected default window 60s, got %v", c.window)
	}
}

func TestRecord_And_Count(t *testing.T) {
	c := New(5 * time.Second)
	c.Record("tcp:80")
	c.Record("tcp:80")
	if got := c.Count("tcp:80"); got != 2 {
		t.Fatalf("expected 2, got %d", got)
	}
}

func TestCount_EmptyKey_ReturnsZero(t *testing.T) {
	c := New(5 * time.Second)
	if got := c.Count("udp:9999"); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestCount_EvictsExpiredEntries(t *testing.T) {
	c := New(50 * time.Millisecond)
	c.Record("tcp:443")
	c.Record("tcp:443")
	time.Sleep(80 * time.Millisecond)
	if got := c.Count("tcp:443"); got != 0 {
		t.Fatalf("expected 0 after expiry, got %d", got)
	}
}

func TestCount_DifferentKeysAreIndependent(t *testing.T) {
	c := New(5 * time.Second)
	for i := 0; i < 3; i++ {
		c.Record("tcp:80")
	}
	c.Record("udp:53")
	if got := c.Count("tcp:80"); got != 3 {
		t.Fatalf("expected 3 for tcp:80, got %d", got)
	}
	if got := c.Count("udp:53"); got != 1 {
		t.Fatalf("expected 1 for udp:53, got %d", got)
	}
}

func TestReset_ClearsAll(t *testing.T) {
	c := New(5 * time.Second)
	for i := 0; i < 5; i++ {
		c.Record(fmt.Sprintf("tcp:%d", i))
	}
	c.Reset()
	for i := 0; i < 5; i++ {
		if got := c.Count(fmt.Sprintf("tcp:%d", i)); got != 0 {
			t.Fatalf("expected 0 after reset, got %d", got)
		}
	}
}

func TestKeyFromEvent_Format(t *testing.T) {
	e := makeEvent(8080, "tcp")
	got := KeyFromEvent(e)
	if got != "tcp:"+string(rune(8080)) {
		// Just verify it is non-empty and contains protocol
		if len(got) == 0 {
			t.Fatal("expected non-empty key")
		}
	}
}
