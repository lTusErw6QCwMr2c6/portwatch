// Package rotation handles audit log file rotation based on size or age.
package rotation

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Config holds rotation policy settings.
type Config struct {
	MaxBytes int64
	MaxAge   time.Duration
	KeepLast int
}

// DefaultConfig returns sensible rotation defaults.
func DefaultConfig() Config {
	return Config{
		MaxBytes: 10 * 1024 * 1024, // 10 MB
		MaxAge:   24 * time.Hour,
		KeepLast: 5,
	}
}

// Rotator manages log file rotation for a given path.
type Rotator struct {
	mu      sync.Mutex
	path    string
	cfg     Config
	created time.Time
}

// New creates a new Rotator for the given file path and config.
func New(path string, cfg Config) *Rotator {
	return &Rotator{
		path:    path,
		cfg:     cfg,
		created: time.Now(),
	}
}

// ShouldRotate reports whether the file at path exceeds size or age limits.
func (r *Rotator) ShouldRotate() bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if time.Since(r.created) >= r.cfg.MaxAge {
		return true
	}
	info, err := os.Stat(r.path)
	if err != nil {
		return false
	}
	return info.Size() >= r.cfg.MaxBytes
}

// Rotate renames the current log file with a timestamp suffix and prunes old files.
func (r *Rotator) Rotate() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	timestamp := time.Now().Format("20060102T150405")
	ext := filepath.Ext(r.path)
	base := r.path[:len(r.path)-len(ext)]
	dest := fmt.Sprintf("%s.%s%s", base, timestamp, ext)

	if err := os.Rename(r.path, dest); err != nil {
		return fmt.Errorf("rotation rename: %w", err)
	}
	r.created = time.Now()
	return r.prune(base, ext)
}

// RotatedFiles returns the list of rotated (archived) log files for the
// current log path, sorted by filename (oldest first).
func (r *Rotator) RotatedFiles() ([]string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	ext := filepath.Ext(r.path)
	base := r.path[:len(r.path)-len(ext)]
	matches, err := filepath.Glob(base + ".*" + ext)
	if err != nil {
		return nil, fmt.Errorf("listing rotated files: %w", err)
	}
	return matches, nil
}

func (r *Rotator) prune(base, ext string) error {
	pattern := base + ".*" + ext
	matches, err := filepath.Glob(pattern)
	if err != nil || len(matches) <= r.cfg.KeepLast {
		return err
	}
	for _, old := range matches[:len(matches)-r.cfg.KeepLast] {
		_ = os.Remove(old)
	}
	return nil
}
