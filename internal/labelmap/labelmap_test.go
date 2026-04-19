package labelmap_test

import (
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/labelmap"
	"github.com/user/portwatch/internal/scanner"
)

func makeEvent(port int, proto string) alert.Event {
	return alert.Event{
		Port: scanner.Port{Number: port, Protocol: proto},
		Type: alert.Opened,
	}
}

func TestApply_NoRules_ReturnsNil(t *testing.T) {
	m := labelmap.New(nil)
	labels := m.Apply(makeEvent(80, "tcp"))
	if labels != nil {
		t.Fatalf("expected nil, got %v", labels)
	}
}

func TestApply_MatchingRule_ReturnsLabels(t *testing.T) {
	m := labelmap.New([]labelmap.Rule{
		{Port: 80, Protocol: "tcp", Labels: []string{"http", "web"}},
	})
	labels := m.Apply(makeEvent(80, "tcp"))
	if len(labels) != 2 || labels[0] != "http" || labels[1] != "web" {
		t.Fatalf("unexpected labels: %v", labels)
	}
}

func TestApply_ProtocolMismatch_ReturnsNil(t *testing.T) {
	m := labelmap.New([]labelmap.Rule{
		{Port: 80, Protocol: "udp", Labels: []string{"dns"}},
	})
	labels := m.Apply(makeEvent(80, "tcp"))
	if len(labels) != 0 {
		t.Fatalf("expected no labels, got %v", labels)
	}
}

func TestApply_EmptyProtocol_MatchesAny(t *testing.T) {
	m := labelmap.New([]labelmap.Rule{
		{Port: 443, Protocol: "", Labels: []string{"tls"}},
	})
	if l := m.Apply(makeEvent(443, "tcp")); len(l) != 1 {
		t.Fatalf("expected 1 label for tcp, got %v", l)
	}
	if l := m.Apply(makeEvent(443, "udp")); len(l) != 1 {
		t.Fatalf("expected 1 label for udp, got %v", l)
	}
}

func TestApply_DeduplicatesLabels(t *testing.T) {
	m := labelmap.New([]labelmap.Rule{
		{Port: 22, Labels: []string{"ssh"}},
		{Port: 22, Labels: []string{"ssh", "secure"}},
	})
	labels := m.Apply(makeEvent(22, "tcp"))
	count := 0
	for _, l := range labels {
		if l == "ssh" {
			count++
		}
	}
	if count != 1 {
		t.Fatalf("expected ssh deduplicated, got %v", labels)
	}
}

func TestReset_ClearsRules(t *testing.T) {
	m := labelmap.New([]labelmap.Rule{
		{Port: 80, Labels: []string{"web"}},
	})
	m.Reset()
	if l := m.Apply(makeEvent(80, "tcp")); len(l) != 0 {
		t.Fatalf("expected no labels after reset, got %v", l)
	}
}
