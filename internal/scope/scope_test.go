package scope

import (
	"sort"
	"testing"
)

func TestNew_Empty(t *testing.T) {
	s := New()
	if s.Len() != 0 {
		t.Fatalf("expected 0, got %d", s.Len())
	}
}

func TestAdd_And_Match(t *testing.T) {
	s := New()
	s.Add(Entry{Name: "http", Ports: []uint16{80, 8080}, Protocol: "tcp"})

	name, ok := s.Match(80, "tcp")
	if !ok || name != "http" {
		t.Fatalf("expected http/true, got %s/%v", name, ok)
	}
}

func TestMatch_ProtocolMismatch_ReturnsFalse(t *testing.T) {
	s := New()
	s.Add(Entry{Name: "dns", Ports: []uint16{53}, Protocol: "udp"})

	_, ok := s.Match(53, "tcp")
	if ok {
		t.Fatal("expected no match for protocol mismatch")
	}
}

func TestMatch_EmptyProtocol_MatchesBoth(t *testing.T) {
	s := New()
	s.Add(Entry{Name: "any53", Ports: []uint16{53}, Protocol: ""})

	for _, proto := range []string{"tcp", "udp"} {
		_, ok := s.Match(53, proto)
		if !ok {
			t.Fatalf("expected match for protocol %s", proto)
		}
	}
}

func TestRemove_DeletesEntry(t *testing.T) {
	s := New()
	s.Add(Entry{Name: "ssh", Ports: []uint16{22}, Protocol: "tcp"})
	s.Remove("ssh")

	_, ok := s.Match(22, "tcp")
	if ok {
		t.Fatal("expected no match after removal")
	}
}

func TestNames_ReturnsAll(t *testing.T) {
	s := New()
	s.Add(Entry{Name: "http", Ports: []uint16{80}})
	s.Add(Entry{Name: "https", Ports: []uint16{443}})

	names := s.Names()
	sort.Strings(names)
	if len(names) != 2 || names[0] != "http" || names[1] != "https" {
		t.Fatalf("unexpected names: %v", names)
	}
}

func TestMatch_UnknownPort_ReturnsFalse(t *testing.T) {
	s := New()
	s.Add(Entry{Name: "http", Ports: []uint16{80}, Protocol: "tcp"})

	_, ok := s.Match(9999, "tcp")
	if ok {
		t.Fatal("expected no match for unknown port")
	}
}
