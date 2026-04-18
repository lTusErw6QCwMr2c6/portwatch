package prescan_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/prescan"
)

func TestDefaultConfig(t *testing.T) {
	cfg := prescan.DefaultConfig()
	if cfg.Retries <= 0 {
		t.Fatalf("expected positive Retries, got %d", cfg.Retries)
	}
	if cfg.PortStart <= 0 {
		t.Fatalf("expected positive PortStart, got %d", cfg.PortStart)
	}
	if cfg.PortEnd <= cfg.PortStart {
		t.Fatalf("PortEnd %d must exceed PortStart %d", cfg.PortEnd, cfg.PortStart)
	}
	if cfg.Delay <= 0 {
		t.Fatal("expected positive Delay")
	}
}

func TestRun_InvalidRange_ReturnsError(t *testing.T) {
	cfg := prescan.Config{
		PortStart: 9999,
		PortEnd:   1,
		Retries:   1,
		Delay:     time.Millisecond,
	}
	_, err := prescan.Run(cfg)
	if err == nil {
		t.Fatal("expected error for invalid port range, got nil")
	}
}

func TestRun_ValidRange_ReturnsResult(t *testing.T) {
	cfg := prescan.Config{
		PortStart: 1,
		PortEnd:   1024,
		Retries:   2,
		Delay:     time.Millisecond,
	}
	res, err := prescan.Run(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.ScannedAt.IsZero() {
		t.Error("expected non-zero ScannedAt")
	}
	if res.Duration <= 0 {
		t.Error("expected positive Duration")
	}
}

func TestRun_RetriesExhausted_ReturnsError(t *testing.T) {
	// Force failure by using an invalid range with zero retries handled via 1 retry.
	cfg := prescan.Config{
		PortStart: -1,
		PortEnd:   -1,
		Retries:   2,
		Delay:     time.Millisecond,
	}
	_, err := prescan.Run(cfg)
	if err == nil {
		t.Fatal("expected error after exhausted retries")
	}
}
