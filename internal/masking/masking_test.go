package masking

import (
	"strings"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Level != LevelPartial {
		t.Fatalf("expected LevelPartial, got %v", cfg.Level)
	}
	if cfg.Placeholder == "" {
		t.Fatal("expected non-empty placeholder")
	}
	if len(cfg.Fields) == 0 {
		t.Fatal("expected default fields")
	}
}

func TestApply_None_ReturnsOriginal(t *testing.T) {
	m := New(Config{Level: LevelNone, Placeholder: "[M]"})
	if got := m.Apply("secret"); got != "secret" {
		t.Fatalf("expected 'secret', got %q", got)
	}
}

func TestApply_Full_ReturnsPlaceholder(t *testing.T) {
	m := New(Config{Level: LevelFull, Placeholder: "[MASKED]"})
	if got := m.Apply("nginx"); got != "[MASKED]" {
		t.Fatalf("expected placeholder, got %q", got)
	}
}

func TestApply_Partial_MasksMiddle(t *testing.T) {
	m := New(Config{Level: LevelPartial, Placeholder: "[M]"})
	got := m.Apply("nginx")
	if !strings.HasPrefix(got, "n") || !strings.HasSuffix(got, "x") {
		t.Fatalf("unexpected partial mask: %q", got)
	}
	if !strings.Contains(got, "*") {
		t.Fatalf("expected asterisks in partial mask: %q", got)
	}
}

func TestApply_ShortString(t *testing.T) {
	m := New(Config{Level: LevelPartial})
	got := m.Apply("ab")
	if got != "**" {
		t.Fatalf("expected '**', got %q", got)
	}
}

func TestApplyField_MasksWhenInList(t *testing.T) {
	m := New(Config{Level: LevelFull, Placeholder: "[X]", Fields: []string{"process"}})
	if got := m.ApplyField("process", "nginx"); got != "[X]" {
		t.Fatalf("expected placeholder, got %q", got)
	}
}

func TestApplyField_SkipsWhenNotInList(t *testing.T) {
	m := New(Config{Level: LevelFull, Placeholder: "[X]", Fields: []string{"process"}})
	if got := m.ApplyField("port", "8080"); got != "8080" {
		t.Fatalf("expected original value, got %q", got)
	}
}

func TestShouldMask_CaseInsensitive(t *testing.T) {
	m := New(Config{Fields: []string{"Process"}})
	if !m.ShouldMask("process") {
		t.Fatal("expected ShouldMask to return true for 'process'")
	}
}

func TestApply_EmptyString_ReturnsEmpty(t *testing.T) {
	m := New(Config{Level: LevelFull, Placeholder: "[M]"})
	if got := m.Apply(""); got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
}
