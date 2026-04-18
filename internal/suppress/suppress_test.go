package suppress

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func makeEvent(proto, addr string) alert.Event {
	return alert.Event{
		Port: scanner.Port{Protocol: proto, Address: addr},
	}
}

func TestAllow_NotSuppressed_ReturnsTrue(t *testing.T) {
	s := New()
	ev := makeEvent("tcp", "0.0.0.0:8080")
	if !s.Allow(ev) {
		t.Fatal("expected event to be allowed")
	}
}

func TestSuppress_BlocksEvent(t *testing.T) {
	s := New()
	s.Suppress("tcp:0.0.0.0:9000", 10*time.Minute, "maintenance")
	suppressed, reason := s.IsSuppressed("tcp:0.0.0.0:9000")
	if !suppressed {
		t.Fatal("expected key to be suppressed")
	}
	if reason != "maintenance" {
		t.Fatalf("expected reason 'maintenance', got %q", reason)
	}
}

func TestSuppress_Expired_AllowsEvent(t *testing.T) {
	s := New()
	fixed := time.Now()
	s.now = func() time.Time { return fixed }
	s.Suppress("tcp:127.0.0.1:443", 1*time.Millisecond, "test")
	s.now = func() time.Time { return fixed.Add(1 * time.Second) }
	suppressed, _ := s.IsSuppressed("tcp:127.0.0.1:443")
	if suppressed {
		t.Fatal("expected suppression to have expired")
	}
}

func TestRemove_LiftsSuppression(t *testing.T) {
	s := New()
	s.Suppress("udp:0.0.0.0:53", 1*time.Hour, "dns")
	s.Remove("udp:0.0.0.0:53")
	suppressed, _ := s.IsSuppressed("udp:0.0.0.0:53")
	if suppressed {
		t.Fatal("expected suppression to be lifted after Remove")
	}
}

func TestReset_ClearsAll(t *testing.T) {
	s := New()
	s.Suppress("tcp:0.0.0.0:80", 1*time.Hour, "a")
	s.Suppress("tcp:0.0.0.0:443", 1*time.Hour, "b")
	s.Reset()
	for _, key := range []string{"tcp:0.0.0.0:80", "tcp:0.0.0.0:443"} {
		if ok, _ := s.IsSuppressed(key); ok {
			t.Fatalf("expected key %q to be cleared after Reset", key)
		}
	}
}

func TestAllow_UsesProtoAndAddress(t *testing.T) {
	s := New()
	s.Suppress("tcp:0.0.0.0:8080", 1*time.Hour, "blocked")
	ev := makeEvent("tcp", "0.0.0.0:8080")
	if s.Allow(ev) {
		t.Fatal("expected event to be suppressed")
	}
}
