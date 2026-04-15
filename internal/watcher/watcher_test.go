package watcher

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/logger"
)

func newTestLogger(buf *bytes.Buffer) *logger.Logger {
	return logger.New(buf, logger.LevelInfo)
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Interval != 5*time.Second {
		t.Errorf("expected 5s interval, got %v", cfg.Interval)
	}
	if cfg.MinPort != 1 {
		t.Errorf("expected MinPort 1, got %d", cfg.MinPort)
	}
	if cfg.MaxPort != 65535 {
		t.Errorf("expected MaxPort 65535, got %d", cfg.MaxPort)
	}
	if len(cfg.Protocols) == 0 {
		t.Error("expected at least one protocol")
	}
}

func TestNew_ReturnsWatcher(t *testing.T) {
	var buf bytes.Buffer
	log := newTestLogger(&buf)
	cfg := DefaultConfig()

	w := New(cfg, log)
	if w == nil {
		t.Fatal("expected non-nil watcher")
	}
	if w.stop == nil {
		t.Error("expected stop channel to be initialised")
	}
}

func TestStop_StopsWatcher(t *testing.T) {
	var buf bytes.Buffer
	log := newTestLogger(&buf)

	cfg := DefaultConfig()
	cfg.Interval = 50 * time.Millisecond
	// Use a narrow port range to keep the scan fast in tests.
	cfg.MinPort = 1
	cfg.MaxPort = 1024

	w := New(cfg, log)

	done := make(chan struct{})
	go func() {
		w.Start()
		close(done)
	}()

	// Give the watcher time to perform at least one tick.
	time.Sleep(120 * time.Millisecond)
	w.Stop()

	select {
	case <-done:
		// success
	case <-time.After(2 * time.Second):
		t.Error("watcher did not stop within timeout")
	}

	output := buf.String()
	if !strings.Contains(output, "portwatch started") {
		t.Errorf("expected start log, got: %s", output)
	}
	if !strings.Contains(output, "portwatch stopped") {
		t.Errorf("expected stop log, got: %s", output)
	}
}
