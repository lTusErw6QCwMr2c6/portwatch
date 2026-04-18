package correlation

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func makeEvent(port int) alert.Event {
	return alert.Event{
		Kind: alert.Opened,
		Port: scanner.Port{Number: port, Proto: "tcp"},
		Time: time.Now(),
	}
}

func TestNew_NotNil(t *testing.T) {
	c := New(nil)
	if c == nil {
		t.Fatal("expected non-nil correlator")
	}
}

func TestEvaluate_NoMatch_BelowThreshold(t *testing.T) {
	rules := []Rule{{Name: "burst", MinPorts: 5, Window: time.Minute}}
	c := New(rules)
	for i := 0; i < 3; i++ {
		c.Add(makeEvent(8000 + i))
	}
	matches := c.Evaluate()
	if len(matches) != 0 {
		t.Fatalf("expected 0 matches, got %d", len(matches))
	}
}

func TestEvaluate_Match_AtThreshold(t *testing.T) {
	rules := []Rule{{Name: "burst", MinPorts: 3, Window: time.Minute}}
	c := New(rules)
	for i := 0; i < 3; i++ {
		c.Add(makeEvent(9000 + i))
	}
	matches := c.Evaluate()
	if len(matches) != 1 {
		t.Fatalf("expected 1 match, got %d", len(matches))
	}
	if matches[0].Rule != "burst" {
		t.Errorf("unexpected rule name: %s", matches[0].Rule)
	}
}

func TestReset_ClearsBucket(t *testing.T) {
	rules := []Rule{{Name: "burst", MinPorts: 1, Window: time.Minute}}
	c := New(rules)
	c.Add(makeEvent(1234))
	c.Reset()
	matches := c.Evaluate()
	if len(matches) != 0 {
		t.Fatalf("expected 0 matches after reset, got %d", len(matches))
	}
}

func TestMatch_String(t *testing.T) {
	m := Match{
		Rule:       "sweep",
		Events:     []alert.Event{makeEvent(80)},
		DetectedAt: time.Now(),
	}
	s := m.String()
	if s == "" {
		t.Error("expected non-empty string")
	}
}

func TestEvaluate_MultipleRules(t *testing.T) {
	rules := []Rule{
		{Name: "small", MinPorts: 2, Window: time.Minute},
		{Name: "large", MinPorts: 10, Window: time.Minute},
	}
	c := New(rules)
	for i := 0; i < 4; i++ {
		c.Add(makeEvent(7000 + i))
	}
	matches := c.Evaluate()
	if len(matches) != 1 {
		t.Fatalf("expected 1 match, got %d", len(matches))
	}
	if matches[0].Rule != "small" {
		t.Errorf("expected 'small' rule, got %s", matches[0].Rule)
	}
}
