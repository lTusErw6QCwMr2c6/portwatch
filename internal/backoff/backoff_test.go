package backoff

import (
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Strategy != Exponential {
		t.Errorf("expected Exponential, got %v", cfg.Strategy)
	}
	if cfg.MaxRetries != 5 {
		t.Errorf("expected 5 retries, got %d", cfg.MaxRetries)
	}
}

func TestNext_ExponentialGrowth(t *testing.T) {
	cfg := Config{
		Strategy:   Exponential,
		BaseDelay:  100 * time.Millisecond,
		MaxDelay:   10 * time.Second,
		Multiplier: 2.0,
		MaxRetries: 4,
	}
	b := New(cfg)
	expected := []time.Duration{
		100 * time.Millisecond,
		200 * time.Millisecond,
		400 * time.Millisecond,
		800 * time.Millisecond,
	}
	for i, want := range expected {
		got, ok := b.Next()
		if !ok {
			t.Fatalf("attempt %d: expected ok=true", i)
		}
		if got != want {
			t.Errorf("attempt %d: got %v, want %v", i, got, want)
		}
	}
	_, ok := b.Next()
	if ok {
		t.Error("expected ok=false after max retries")
	}
}

func TestNext_FixedStrategy(t *testing.T) {
	cfg := Config{
		Strategy:   Fixed,
		BaseDelay:  250 * time.Millisecond,
		MaxDelay:   5 * time.Second,
		Multiplier: 2.0,
		MaxRetries: 3,
	}
	b := New(cfg)
	for i := 0; i < 3; i++ {
		got, ok := b.Next()
		if !ok {
			t.Fatalf("attempt %d: expected ok", i)
		}
		if got != 250*time.Millisecond {
			t.Errorf("attempt %d: got %v, want 250ms", i, got)
		}
	}
}

func TestNext_CapsAtMaxDelay(t *testing.T) {
	cfg := Config{
		Strategy:   Exponential,
		BaseDelay:  1 * time.Second,
		MaxDelay:   2 * time.Second,
		Multiplier: 10.0,
		MaxRetries: 3,
	}
	b := New(cfg)
	b.Next() // 1s
	got, _ := b.Next() // would be 10s, capped at 2s
	if got != 2*time.Second {
		t.Errorf("expected delay capped at 2s, got %v", got)
	}
}

func TestReset_ClearsAttempt(t *testing.T) {
	b := New(DefaultConfig())
	b.Next()
	b.Next()
	if b.Attempt() != 2 {
		t.Errorf("expected attempt=2, got %d", b.Attempt())
	}
	b.Reset()
	if b.Attempt() != 0 {
		t.Errorf("expected attempt=0 after reset, got %d", b.Attempt())
	}
}
