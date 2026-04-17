package rotation_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/portwatch/internal/rotation"
)

func tempLog(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "audit*.log")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.WriteString(content)
	_ = f.Close()
	return f.Name()
}

func TestDefaultConfig(t *testing.T) {
	cfg := rotation.DefaultConfig()
	if cfg.MaxBytes <= 0 {
		t.Error("expected positive MaxBytes")
	}
	if cfg.MaxAge <= 0 {
		t.Error("expected positive MaxAge")
	}
	if cfg.KeepLast <= 0 {
		t.Error("expected positive KeepLast")
	}
}

func TestShouldRotate_FalseForSmallNewFile(t *testing.T) {
	path := tempLog(t, "hello")
	cfg := rotation.DefaultConfig()
	r := rotation.New(path, cfg)
	if r.ShouldRotate() {
		t.Error("expected no rotation needed for small new file")
	}
}

func TestShouldRotate_TrueWhenExceedsSize(t *testing.T) {
	path := tempLog(t, "data")
	cfg := rotation.DefaultConfig()
	cfg.MaxBytes = 2 // very small threshold
	r := rotation.New(path, cfg)
	if !r.ShouldRotate() {
		t.Error("expected rotation needed when file exceeds MaxBytes")
	}
}

func TestShouldRotate_TrueWhenAgeExceeded(t *testing.T) {
	path := tempLog(t, "x")
	cfg := rotation.DefaultConfig()
	cfg.MaxAge = time.Nanosecond
	time.Sleep(2 * time.Millisecond)
	r := rotation.New(path, cfg)
	if !r.ShouldRotate() {
		t.Error("expected rotation needed when age exceeded")
	}
}

func TestRotate_RenamesFile(t *testing.T) {
	path := tempLog(t, "log data")
	cfg := rotation.DefaultConfig()
	r := rotation.New(path, cfg)
	if err := r.Rotate(); err != nil {
		t.Fatalf("Rotate() error: %v", err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("expected original file to be renamed away")
	}
	matches, _ := filepath.Glob(path[:len(path)-len(filepath.Ext(path))] + ".*" + filepath.Ext(path))
	if len(matches) == 0 {
		t.Error("expected at least one rotated file")
	}
}
