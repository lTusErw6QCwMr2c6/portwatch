package eventgrep

import (
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func makeEvent(port int, proto, process string, t alert.EventType) alert.Event {
	return alert.Event{
		Type: t,
		Port: scanner.Port{
			Port:     port,
			Protocol: proto,
			Process:  process,
		},
	}
}

func TestNew_NotNil(t *testing.T) {
	g := New()
	if g == nil {
		t.Fatal("expected non-nil Grepper")
	}
	if g.Len() != 0 {
		t.Fatalf("expected 0 matchers, got %d", g.Len())
	}
}

func TestAdd_InvalidPattern_ReturnsError(t *testing.T) {
	g := New()
	if err := g.Add("["); err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}

func TestMatch_NoRules_ReturnsFalse(t *testing.T) {
	g := New()
	e := makeEvent(8080, "tcp", "nginx", alert.EventOpened)
	if g.Match(e) {
		t.Fatal("expected no match with no rules")
	}
}

func TestMatch_ProtocolPattern_ReturnsTrue(t *testing.T) {
	g := New()
	_ = g.Add("tcp", "protocol")
	e := makeEvent(443, "tcp", "sshd", alert.EventOpened)
	if !g.Match(e) {
		t.Fatal("expected match on protocol")
	}
}

func TestMatch_ProcessPattern_ReturnsFalse_WhenMismatch(t *testing.T) {
	g := New()
	_ = g.Add("nginx", "process")
	e := makeEvent(22, "tcp", "sshd", alert.EventOpened)
	if g.Match(e) {
		t.Fatal("expected no match")
	}
}

func TestFilter_ReturnsMatchingSubset(t *testing.T) {
	g := New()
	_ = g.Add("opened", "type")

	events := []alert.Event{
		makeEvent(80, "tcp", "nginx", alert.EventOpened),
		makeEvent(443, "tcp", "nginx", alert.EventClosed),
		makeEvent(8080, "tcp", "caddy", alert.EventOpened),
	}

	result := g.Filter(events)
	if len(result) != 2 {
		t.Fatalf("expected 2 results, got %d", len(result))
	}
}

func TestMatch_DefaultFields_SearchesAll(t *testing.T) {
	g := New()
	_ = g.Add("udp") // no explicit fields — uses all defaults

	e := makeEvent(53, "udp", "dnsmasq", alert.EventOpened)
	if !g.Match(e) {
		t.Fatal("expected match on default fields")
	}
}
