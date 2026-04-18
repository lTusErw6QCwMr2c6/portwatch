package policy_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/user/portwatch/internal/policy"
)

func writePolicyFile(t *testing.T, rules []map[string]interface{}) string {
	t.Helper()
	f, err := os.CreateTemp("", "policy-*.json")
	if err != nil {
		t.Fatal(err)
	}
	if err := json.NewEncoder(f).Encode(rules); err != nil {
		t.Fatal(err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestLoad_ValidFile(t *testing.T) {
	path := writePolicyFile(t, []map[string]interface{}{
		{"name": "deny-22", "port": 22, "protocol": "tcp", "action": "deny"},
	})
	p, err := policy.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == nil {
		t.Fatal("expected non-nil policy")
	}
}

func TestLoad_UnknownAction_ReturnsError(t *testing.T) {
	path := writePolicyFile(t, []map[string]interface{}{
		{"name": "bad", "port": 80, "action": "explode"},
	})
	_, err := policy.Load(path)
	if err == nil {
		t.Fatal("expected error for unknown action")
	}
}

func TestLoad_MissingFile_ReturnsError(t *testing.T) {
	_, err := policy.Load("/nonexistent/policy.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoad_AppliesRuleCorrectly(t *testing.T) {
	path := writePolicyFile(t, []map[string]interface{}{
		{"name": "warn-443", "port": 443, "protocol": "tcp", "action": "warn"},
	})
	p, err := policy.Load(path)
	if err != nil {
		t.Fatal(err)
	}
	action, name := p.Evaluate(makeEvent(443, "tcp"))
	if action != policy.ActionWarn || name != "warn-443" {
		t.Errorf("got %s/%s", action, name)
	}
}
