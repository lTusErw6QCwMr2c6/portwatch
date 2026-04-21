package eventscope

import (
	"testing"

	"github.com/user/portwatch/internal/alert"
)

func makeEvent(port int, proto string) alert.Event {
	return alert.Event{Port: port, Protocol: proto, Type: alert.Opened}
}

func TestNew_Empty(t *testing.T) {
	s := New()
	if s.Len() != 0 {
		t.Fatalf("expected 0 scopes, got %d", s.Len())
	}
}

func TestAdd_ValidScope(t *testing.T) {
	s := New()
	err := s.Add(Scope{Name: "web", Protocol: "tcp", PortMin: 80, PortMax: 443})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Len() != 1 {
		t.Fatalf("expected 1 scope, got %d", s.Len())
	}
}

func TestAdd_DuplicateName_ReturnsError(t *testing.T) {
	s := New()
	sc := Scope{Name: "web", Protocol: "tcp", PortMin: 80, PortMax: 443}
	_ = s.Add(sc)
	err := s.Add(sc)
	if err == nil {
		t.Fatal("expected error for duplicate scope name")
	}
}

func TestAdd_InvalidRange_ReturnsError(t *testing.T) {
	s := New()
	err := s.Add(Scope{Name: "bad", Protocol: "tcp", PortMin: 500, PortMax: 100})
	if err == nil {
		t.Fatal("expected error for invalid port range")
	}
}

func TestAdd_EmptyName_ReturnsError(t *testing.T) {
	s := New()
	err := s.Add(Scope{Name: "", Protocol: "tcp", PortMin: 1, PortMax: 100})
	if err == nil {
		t.Fatal("expected error for empty scope name")
	}
}

func TestMatch_PortInRange(t *testing.T) {
	s := New()
	_ = s.Add(Scope{Name: "web", Protocol: "tcp", PortMin: 80, PortMax: 443})
	matched := s.Match(makeEvent(443, "tcp"))
	if len(matched) != 1 || matched[0] != "web" {
		t.Fatalf("expected [web], got %v", matched)
	}
}

func TestMatch_PortOutsideRange_ReturnsEmpty(t *testing.T) {
	s := New()
	_ = s.Add(Scope{Name: "web", Protocol: "tcp", PortMin: 80, PortMax: 443})
	matched := s.Match(makeEvent(8080, "tcp"))
	if len(matched) != 0 {
		t.Fatalf("expected no match, got %v", matched)
	}
}

func TestMatch_ProtocolMismatch_ReturnsEmpty(t *testing.T) {
	s := New()
	_ = s.Add(Scope{Name: "dns", Protocol: "udp", PortMin: 53, PortMax: 53})
	matched := s.Match(makeEvent(53, "tcp"))
	if len(matched) != 0 {
		t.Fatalf("expected no match on protocol mismatch, got %v", matched)
	}
}

func TestMatch_EmptyProtocol_MatchesAny(t *testing.T) {
	s := New()
	_ = s.Add(Scope{Name: "any53", Protocol: "", PortMin: 53, PortMax: 53})
	for _, proto := range []string{"tcp", "udp"} {
		matched := s.Match(makeEvent(53, proto))
		if len(matched) != 1 {
			t.Fatalf("expected match for protocol %s, got %v", proto, matched)
		}
	}
}

func TestRemove_DeletesScope(t *testing.T) {
	s := New()
	_ = s.Add(Scope{Name: "web", Protocol: "tcp", PortMin: 80, PortMax: 443})
	s.Remove("web")
	if s.Len() != 0 {
		t.Fatalf("expected 0 scopes after remove, got %d", s.Len())
	}
}
