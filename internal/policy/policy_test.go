package policy_test

import (
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/policy"
	"github.com/user/portwatch/internal/scanner"
)

func makeEvent(port uint16, proto string) alert.Event {
	return alert.Event{
		Kind: alert.Opened,
		Port: scanner.Port{Port: port, Protocol: proto},
	}
}

func TestEvaluate_NoRules_DefaultsAllow(t *testing.T) {
	p := policy.New(nil)
	action, name := p.Evaluate(makeEvent(80, "tcp"))
	if action != policy.ActionAllow {
		t.Errorf("expected allow, got %s", action)
	}
	if name != "" {
		t.Errorf("expected empty rule name, got %s", name)
	}
}

func TestEvaluate_MatchingPortRule(t *testing.T) {
	rules := []policy.Rule{
		{Name: "block-22", Port: 22, Protocol: "tcp", Action: policy.ActionDeny},
	}
	p := policy.New(rules)
	action, name := p.Evaluate(makeEvent(22, "tcp"))
	if action != policy.ActionDeny {
		t.Errorf("expected deny, got %s", action)
	}
	if name != "block-22" {
		t.Errorf("expected block-22, got %s", name)
	}
}

func TestEvaluate_ProtocolMismatch_NoMatch(t *testing.T) {
	rules := []policy.Rule{
		{Name: "block-22-tcp", Port: 22, Protocol: "tcp", Action: policy.ActionDeny},
	}
	p := policy.New(rules)
	action, _ := p.Evaluate(makeEvent(22, "udp"))
	if action != policy.ActionAllow {
		t.Errorf("expected allow for udp, got %s", action)
	}
}

func TestEvaluate_WarnAction(t *testing.T) {
	rules := []policy.Rule{
		{Name: "warn-443", Port: 443, Action: policy.ActionWarn},
	}
	p := policy.New(rules)
	action, name := p.Evaluate(makeEvent(443, "tcp"))
	if action != policy.ActionWarn {
		t.Errorf("expected warn, got %s", action)
	}
	if name != "warn-443" {
		t.Errorf("unexpected name: %s", name)
	}
}

func TestEvaluate_FirstMatchWins(t *testing.T) {
	rules := []policy.Rule{
		{Name: "first", Port: 80, Action: policy.ActionDeny},
		{Name: "second", Port: 80, Action: policy.ActionWarn},
	}
	p := policy.New(rules)
	action, name := p.Evaluate(makeEvent(80, "tcp"))
	if action != policy.ActionDeny || name != "first" {
		t.Errorf("expected first deny rule, got %s/%s", action, name)
	}
}

func TestString_ContainsRuleCount(t *testing.T) {
	p := policy.New([]policy.Rule{{Name: "r1"}, {Name: "r2"}})
	s := p.String()
	if s != "Policy(2 rules)" {
		t.Errorf("unexpected string: %s", s)
	}
}
