package eventmatch_test

import (
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/eventmatch"
	"github.com/user/portwatch/internal/scanner"
)

func makeEvent(port int, proto, typ string) alert.Event {
	return alert.Event{
		Port:  scanner.Port{Port: port, Protocol: proto},
		Type:  alert.EventType(typ),
	}
}

func TestMatch_NoRules_ReturnsFalse(t *testing.T) {
	m := eventmatch.New(nil)
	_, ok := m.Match(makeEvent(80, "tcp", "opened"))
	if ok {
		t.Fatal("expected no match")
	}
}

func TestMatch_PortMatch_ReturnsRule(t *testing.T) {
	rules := []eventmatch.Rule{{Port: 80, Label: "http"}}
	m := eventmatch.New(rules)
	r, ok := m.Match(makeEvent(80, "tcp", "opened"))
	if !ok {
		t.Fatal("expected match")
	}
	if r.Label != "http" {
		t.Fatalf("expected label http, got %s", r.Label)
	}
}

func TestMatch_PortMismatch_ReturnsFalse(t *testing.T) {
	rules := []eventmatch.Rule{{Port: 443}}
	m := eventmatch.New(rules)
	_, ok := m.Match(makeEvent(80, "tcp", "opened"))
	if ok {
		t.Fatal("expected no match")
	}
}

func TestMatch_ProtocolMismatch_ReturnsFalse(t *testing.T) {
	rules := []eventmatch.Rule{{Port: 80, Protocol: "udp"}}
	m := eventmatch.New(rules)
	_, ok := m.Match(makeEvent(80, "tcp", "opened"))
	if ok {
		t.Fatal("expected no match on protocol mismatch")
	}
}

func TestMatch_TypeFilter_Matches(t *testing.T) {
	rules := []eventmatch.Rule{{Type: "closed", Label: "gone"}}
	m := eventmatch.New(rules)
	r, ok := m.Match(makeEvent(9000, "tcp", "closed"))
	if !ok {
		t.Fatal("expected match on type")
	}
	if r.Label != "gone" {
		t.Fatalf("unexpected label: %s", r.Label)
	}
}

func TestMatchAll_ReturnsMultiple(t *testing.T) {
	rules := []eventmatch.Rule{
		{Port: 80, Label: "port-rule"},
		{Protocol: "tcp", Label: "proto-rule"},
		{Type: "opened", Label: "type-rule"},
	}
	m := eventmatch.New(rules)
	matched := m.MatchAll(makeEvent(80, "tcp", "opened"))
	if len(matched) != 3 {
		t.Fatalf("expected 3 matches, got %d", len(matched))
	}
}

func TestMatchAll_EmptyInput_ReturnsNil(t *testing.T) {
	m := eventmatch.New([]eventmatch.Rule{{Port: 443}})
	matched := m.MatchAll(makeEvent(80, "tcp", "opened"))
	if len(matched) != 0 {
		t.Fatalf("expected no matches, got %d", len(matched))
	}
}
