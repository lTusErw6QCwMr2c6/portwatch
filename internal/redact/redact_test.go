package redact_test

import (
	"testing"

	"github.com/user/portwatch/internal/redact"
)

func TestApply_MasksPassword(t *testing.T) {
	r := redact.New(redact.DefaultRules())
	input := "connect password=supersecret host=db"
	out := r.Apply(input)
	if contains(out, "supersecret") {
		t.Fatalf("expected password to be redacted, got: %s", out)
	}
	if !contains(out, "password=***") {
		t.Fatalf("expected placeholder, got: %s", out)
	}
}

func TestApply_MasksIP(t *testing.T) {
	r := redact.New(redact.DefaultRules())
	out := r.Apply("connected from 192.168.1.100")
	if contains(out, "192.168.1.100") {
		t.Fatalf("IP should be redacted, got: %s", out)
	}
}

func TestApply_NoMatch_Unchanged(t *testing.T) {
	r := redact.New(redact.DefaultRules())
	input := "nginx worker process"
	if got := r.Apply(input); got != input {
		t.Fatalf("expected unchanged string, got: %s", got)
	}
}

func TestApplyAll_RedactsEach(t *testing.T) {
	r := redact.New(redact.DefaultRules())
	ss := []string{"token=abc123", "safe string"}
	out := r.ApplyAll(ss)
	if contains(out[0], "abc123") {
		t.Fatal("token value should be redacted")
	}
	if out[1] != "safe string" {
		t.Fatal("safe string should be unchanged")
	}
}

func TestContainsSensitive_True(t *testing.T) {
	r := redact.New(redact.DefaultRules())
	if !r.ContainsSensitive("secret=xyz") {
		t.Fatal("expected sensitive match")
	}
}

func TestContainsSensitive_False(t *testing.T) {
	r := redact.New(redact.DefaultRules())
	if r.ContainsSensitive("nginx") {
		t.Fatal("expected no sensitive match")
	}
}

func TestNew_EmptyRules_NoChange(t *testing.T) {
	r := redact.New(nil)
	if got := r.Apply("password=secret"); got != "password=secret" {
		t.Fatalf("expected unchanged with no rules, got: %s", got)
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 ||
		(func() bool {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
			return false
		})())
}
