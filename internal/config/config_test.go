package config

import (
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Interval != 5 {
		t.Errorf("expected interval 5, got %d", cfg.Interval)
	}
	if cfg.PortStart != 1 {
		t.Errorf("expected port-start 1, got %d", cfg.PortStart)
	}
	if cfg.PortEnd != 65535 {
		t.Errorf("expected port-end 65535, got %d", cfg.PortEnd)
	}
	if cfg.LogLevel != "info" {
		t.Errorf("expected log-level 'info', got %q", cfg.LogLevel)
	}
	if cfg.JSON {
		t.Error("expected JSON to be false by default")
	}
}

func TestValidate_ValidConfig(t *testing.T) {
	cfg := DefaultConfig()
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected no error for valid config, got: %v", err)
	}
}

func TestValidate_InvalidInterval(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Interval = 0
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for interval=0, got nil")
	}
}

func TestValidate_PortStartExceedsPortEnd(t *testing.T) {
	cfg := DefaultConfig()
	cfg.PortStart = 9000
	cfg.PortEnd = 1000
	if err := cfg.Validate(); err == nil {
		t.Error("expected error when port-start > port-end, got nil")
	}
}

func TestValidate_InvalidPortRange(t *testing.T) {
	tests := []struct {
		name  string
		start int
		end   int
	}{
		{"start zero", 0, 1000},
		{"end too large", 1, 70000},
		{"start too large", 70000, 80000},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := DefaultConfig()
			cfg.PortStart = tt.start
			cfg.PortEnd = tt.end
			if err := cfg.Validate(); err == nil {
				t.Errorf("expected validation error for start=%d end=%d", tt.start, tt.end)
			}
		})
	}
}

func TestValidate_InvalidLogLevel(t *testing.T) {
	cfg := DefaultConfig()
	cfg.LogLevel = "verbose"
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for invalid log level, got nil")
	}
}

func TestValidate_AllLogLevels(t *testing.T) {
	levels := []string{"debug", "info", "warn", "error"}
	for _, level := range levels {
		t.Run(level, func(t *testing.T) {
			cfg := DefaultConfig()
			cfg.LogLevel = level
			if err := cfg.Validate(); err != nil {
				t.Errorf("expected no error for log-level=%q, got: %v", level, err)
			}
		})
	}
}
