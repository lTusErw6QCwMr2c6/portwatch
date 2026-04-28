package eventlabel

import (
	"bytes"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/logger"
)

func makeEvent(port int, proto string) alert.Event {
	return alert.Event{Port: port, Protocol: proto, Type: alert.Opened}
}

func TestNew_NotNil(t *testing.T) {
	if New() == nil {
		t.Fatal("expected non-nil Labeler")
	}
}

func TestRegister_AddsRule(t *testing.T) {
	l := New()
	if err := l.Register("env", "environment", "prod", nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l.Len() != 1 {
		t.Fatalf("expected 1 rule, got %d", l.Len())
	}
}

func TestRegister_EmptyName_ReturnsError(t *testing.T) {
	l := New()
	if err := l.Register("", "k", "v", nil); err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestRegister_DuplicateName_ReturnsError(t *testing.T) {
	l := New()
	_ = l.Register("r", "k", "v", nil)
	if err := l.Register("r", "k2", "v2", nil); err == nil {
		t.Fatal("expected error for duplicate name")
	}
}

func TestApply_NoRules_ReturnsNil(t *testing.T) {
	l := New()
	if labels := l.Apply(makeEvent(80, "tcp")); labels != nil {
		t.Fatalf("expected nil, got %v", labels)
	}
}

func TestApply_UnconditionalRule_AlwaysMatches(t *testing.T) {
	l := New()
	_ = l.Register("svc", "service", "web", nil)
	labels := l.Apply(makeEvent(80, "tcp"))
	if labels["service"] != "web" {
		t.Fatalf("expected label service=web, got %v", labels)
	}
}

func TestApply_PredicateMismatch_ReturnsNil(t *testing.T) {
	l := New()
	_ = l.Register("high", "tier", "high", func(ev alert.Event) bool {
		return ev.Port > 1024
	})
	if labels := l.Apply(makeEvent(80, "tcp")); labels != nil {
		t.Fatalf("expected nil for port 80, got %v", labels)
	}
}

func TestApply_PredicateMatch_ReturnsLabel(t *testing.T) {
	l := New()
	_ = l.Register("high", "tier", "high", func(ev alert.Event) bool {
		return ev.Port > 1024
	})
	labels := l.Apply(makeEvent(8080, "tcp"))
	if labels["tier"] != "high" {
		t.Fatalf("expected tier=high, got %v", labels)
	}
}

func TestHandler_CallsCallback(t *testing.T) {
	l := New()
	_ = l.Register("env", "env", "staging", nil)
	log := logger.New(bytes.NewBuffer(nil))
	var got map[string]string
	h := NewHandler(l, log, func(_ alert.Event, labels map[string]string) {
		got = labels
	})
	h.Handle(makeEvent(443, "tcp"))
	if got["env"] != "staging" {
		t.Fatalf("expected env=staging, got %v", got)
	}
}
