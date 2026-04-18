package retention

import (
	"os"
	"path/filepath"
	"sort"
	"time"
)

// Config holds retention policy settings.
type Config struct {
	MaxAge     time.Duration
	MaxFiles   int
	GlobPattern string
}

// DefaultConfig returns sensible retention defaults.
func DefaultConfig() Config {
	return Config{
		MaxAge:      7 * 24 * time.Hour,
		MaxFiles:    10,
		GlobPattern: "*.log",
	}
}

// Policy enforces file retention rules on a directory.
type Policy struct {
	dir string
	cfg Config
}

// New creates a new retention Policy.
func New(dir string, cfg Config) *Policy {
	return &Policy{dir: dir, cfg: cfg}
}

// Apply removes files that exceed age or count limits.
func (p *Policy) Apply() ([]string, error) {
	matches, err := filepath.Glob(filepath.Join(p.dir, p.cfg.GlobPattern))
	if err != nil {
		return nil, err
	}

	type entry struct {
		path    string
		modTime time.Time
	}

	var entries []entry
	for _, m := range matches {
		fi, err := os.Stat(m)
		if err != nil {
			continue
		}
		entries = append(entries, entry{path: m, modTime: fi.ModTime()})
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].modTime.After(entries[j].modTime)
	})

	now := time.Now()
	var removed []string
	for i, e := range entries {
		expired := now.Sub(e.modTime) > p.cfg.MaxAge
		excess := i >= p.cfg.MaxFiles
		if expired || excess {
			if err := os.Remove(e.path); err == nil {
				removed = append(removed, e.path)
			}
		}
	}
	return removed, nil
}
