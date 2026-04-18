package retention

import (
	"testing"
	"time"
)

type silentLogger struct{}

func (s *silentLogger) Infof(_ string, _ ...interface{})  {}
func (s *silentLogger) Errorf(_ string, _ ...interface{}) {}

func TestNewRunner_NotNil(t *testing.T) {
	dir := t.TempDir()
	p := New(dir, DefaultConfig())
	r := NewRunner(p, time.Minute, &silentLogger{})
	if r == nil {
		t.Fatal("expected non-nil runner")
	}
}

func TestRunner_StartStop_DoesNotBlock(t *testing.T) {
	dir := t.TempDir()
	p := New(dir, DefaultConfig())
	r := NewRunner(p, time.Hour, &silentLogger{})
	r.Start()
	done := make(chan struct{})
	go func() {
		r.Stop()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("Stop() blocked")
	}
}

func TestRunner_AppliesPolicy(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "stale.log", 48*time.Hour)
	cfg := Config{MaxAge: time.Hour, MaxFiles: 100, GlobPattern: "*.log"}
	p := New(dir, cfg)
	r := NewRunner(p, 20*time.Millisecond, &silentLogger{})
	r.Start()
	time.Sleep(60 * time.Millisecond)
	r.Stop()

	remaining, _ := filepath.Glob(dir + "/*.log")
	if len(remaining) != 0 {
		t.Errorf("expected file removed, got %d remaining", len(remaining))
	}
}
