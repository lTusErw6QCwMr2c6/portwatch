package retention

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeFile(t *testing.T, dir, name string, age time.Duration) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte("data"), 0o644); err != nil {
		t.Fatal(err)
	}
	mt := time.Now().Add(-age)
	if err := os.Chtimes(p, mt, mt); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.MaxFiles != 10 {
		t.Errorf("expected MaxFiles=10, got %d", cfg.MaxFiles)
	}
	if cfg.MaxAge != 7*24*time.Hour {
		t.Errorf("unexpected MaxAge: %v", cfg.MaxAge)
	}
}

func TestApply_RemovesExpiredFiles(t *testing.T) {
	dir := t.TempDir()
	cfg := Config{MaxAge: time.Hour, MaxFiles: 100, GlobPattern: "*.log"}
	p := New(dir, cfg)

	writeFile(t, dir, "old.log", 2*time.Hour)
	writeFile(t, dir, "new.log", 10*time.Minute)

	removed, err := p.Apply()
	if err != nil {
		t.Fatal(err)
	}
	if len(removed) != 1 {
		t.Errorf("expected 1 removed, got %d", len(removed))
	}
}

func TestApply_RemovesExcessFiles(t *testing.T) {
	dir := t.TempDir()
	cfg := Config{MaxAge: 30 * 24 * time.Hour, MaxFiles: 2, GlobPattern: "*.log"}
	p := New(dir, cfg)

	for i, name := range []string{"a.log", "b.log", "c.log"} {
		writeFile(t, dir, name, time.Duration(i)*time.Minute)
	}

	removed, err := p.Apply()
	if err != nil {
		t.Fatal(err)
	}
	if len(removed) != 1 {
		t.Errorf("expected 1 removed, got %d", len(removed))
	}
}

func TestApply_EmptyDir_NoError(t *testing.T) {
	dir := t.TempDir()
	cfg := DefaultConfig()
	p := New(dir, cfg)
	removed, err := p.Apply()
	if err != nil {
		t.Fatal(err)
	}
	if len(removed) != 0 {
		t.Errorf("expected 0 removed, got %d", len(removed))
	}
}
