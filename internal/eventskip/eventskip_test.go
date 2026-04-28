package eventskip_test

import (
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/eventskip"
	"github.com/user/portwatch/internal/scanner"
)

func makeEvent(port int, proto string) alert.Event {
	return alert.Event{
		Port: scanner.Port{Number: port, Protocol: proto},
	}
}

func TestNew_NotNil(t *testing.T) {
	s := eventskip.New()
	if s == nil {
		t.Fatal("expected non-nil Skipper")
	}
}

func TestRegister_AddsRule(t *testing.T) {
	s := eventskip.New()
	err := s.Register("block80", func(e alert.Event) bool { return e.Port.Number == 80 })
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Len() != 1 {
		t.Fatalf("expected 1 rule, got %d", s.Len())
	}
}

func TestRegister_EmptyName_ReturnsError(t *testing.T) {
	s := eventskip.New()
	err := s.Register("", func(e alert.Event) bool { return true })
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestRegister_DuplicateName_ReturnsError(t *testing.T) {
	s := eventskip.New()
	p := func(e alert.Event) bool { return false }
	_ = s.Register("dup", p)
	err := s.Register("dup", p)
	if err == nil {
		t.Fatal("expected error for duplicate name")
	}
}

func TestSkip_MatchingRule_ReturnsTrue(t *testing.T) {
	s := eventskip.New()
	_ = s.Register("block443", func(e alert.Event) bool { return e.Port.Number == 443 })
	e := makeEvent(443, "tcp")
	if !s.Skip(e) {
		t.Error("expected event to be skipped")
	}
}

func TestSkip_NoMatchingRule_ReturnsFalse(t *testing.T) {
	s := eventskip.New()
	_ = s.Register("block443", func(e alert.Event) bool { return e.Port.Number == 443 })
	e := makeEvent(8080, "tcp")
	if s.Skip(e) {
		t.Error("expected event to not be skipped")
	}
}

func TestSkip_NoRules_ReturnsFalse(t *testing.T) {
	s := eventskip.New()
	e := makeEvent(22, "tcp")
	if s.Skip(e) {
		t.Error("expected false with no rules")
	}
}

func TestDeregister_RemovesRule(t *testing.T) {
	s := eventskip.New()
	_ = s.Register("block22", func(e alert.Event) bool { return e.Port.Number == 22 })
	removed := s.Deregister("block22")
	if !removed {
		t.Fatal("expected rule to be removed")
	}
	if s.Len() != 0 {
		t.Fatalf("expected 0 rules after deregister, got %d", s.Len())
	}
	if s.Skip(makeEvent(22, "tcp")) {
		t.Error("expected no skip after deregister")
	}
}

func TestDeregister_UnknownName_ReturnsFalse(t *testing.T) {
	s := eventskip.New()
	if s.Deregister("nonexistent") {
		t.Error("expected false for unknown rule name")
	}
}
