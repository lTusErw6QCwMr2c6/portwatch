package eventclassifier_test

import (
	"testing"

	"github.com/example/portwatch/internal/alert"
	"github.com/example/portwatch/internal/eventclassifier"
	"github.com/example/portwatch/internal/scanner"
)

func makeEvent(port int, proto string) alert.Event {
	return alert.Event{
		Port: scanner.Port{Number: port, Protocol: proto},
		Type: alert.Opened,
	}
}

func TestNew_NotNil(t *testing.T) {
	c := eventclassifier.New()
	if c == nil {
		t.Fatal("expected non-nil classifier")
	}
}

func TestClassify_NoRules_ReturnsUnknown(t *testing.T) {
	c := eventclassifier.New()
	cls := c.Classify(makeEvent(80, "tcp"))
	if cls != eventclassifier.ClassUnknown {
		t.Fatalf("expected ClassUnknown, got %s", cls)
	}
}

func TestClassify_MatchingRule_ReturnsClass(t *testing.T) {
	c := eventclassifier.New()
	c.AddRule(eventclassifier.Rule{
		PortStart: 80,
		PortEnd:   80,
		Protocol:  "tcp",
		Class:     eventclassifier.ClassNormal,
	})
	cls := c.Classify(makeEvent(80, "tcp"))
	if cls != eventclassifier.ClassNormal {
		t.Fatalf("expected ClassNormal, got %s", cls)
	}
}

func TestClassify_ProtocolMismatch_ReturnsUnknown(t *testing.T) {
	c := eventclassifier.New()
	c.AddRule(eventclassifier.Rule{
		PortStart: 80,
		PortEnd:   80,
		Protocol:  "tcp",
		Class:     eventclassifier.ClassNormal,
	})
	cls := c.Classify(makeEvent(80, "udp"))
	if cls != eventclassifier.ClassUnknown {
		t.Fatalf("expected ClassUnknown, got %s", cls)
	}
}

func TestClassify_EmptyProtocol_MatchesAny(t *testing.T) {
	c := eventclassifier.New()
	c.AddRule(eventclassifier.Rule{
		PortStart: 1,
		PortEnd:   1024,
		Protocol:  "",
		Class:     eventclassifier.ClassSuspicious,
	})
	for _, proto := range []string{"tcp", "udp"} {
		cls := c.Classify(makeEvent(443, proto))
		if cls != eventclassifier.ClassSuspicious {
			t.Fatalf("proto %s: expected ClassSuspicious, got %s", proto, cls)
		}
	}
}

func TestClassify_FirstRuleWins(t *testing.T) {
	c := eventclassifier.New()
	c.AddRule(eventclassifier.Rule{PortStart: 22, PortEnd: 22, Protocol: "tcp", Class: eventclassifier.ClassCritical})
	c.AddRule(eventclassifier.Rule{PortStart: 1, PortEnd: 1024, Protocol: "tcp", Class: eventclassifier.ClassNormal})
	cls := c.Classify(makeEvent(22, "tcp"))
	if cls != eventclassifier.ClassCritical {
		t.Fatalf("expected ClassCritical, got %s", cls)
	}
}

func TestClassifyAll_ReturnsMapForEachEvent(t *testing.T) {
	c := eventclassifier.New()
	c.AddRule(eventclassifier.Rule{PortStart: 80, PortEnd: 80, Protocol: "tcp", Class: eventclassifier.ClassNormal})
	events := []alert.Event{makeEvent(80, "tcp"), makeEvent(9999, "tcp")}
	out := c.ClassifyAll(events)
	if out[0] != eventclassifier.ClassNormal {
		t.Fatalf("index 0: expected ClassNormal, got %s", out[0])
	}
	if out[1] != eventclassifier.ClassUnknown {
		t.Fatalf("index 1: expected ClassUnknown, got %s", out[1])
	}
}

func TestReset_ClearsRules(t *testing.T) {
	c := eventclassifier.New()
	c.AddRule(eventclassifier.Rule{PortStart: 80, PortEnd: 80, Protocol: "tcp", Class: eventclassifier.ClassNormal})
	c.Reset()
	cls := c.Classify(makeEvent(80, "tcp"))
	if cls != eventclassifier.ClassUnknown {
		t.Fatalf("after reset expected ClassUnknown, got %s", cls)
	}
}
