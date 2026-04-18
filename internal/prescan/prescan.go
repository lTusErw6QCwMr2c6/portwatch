// Package prescan performs a pre-scan warm-up before the main watch loop,
// establishing an initial port baseline to avoid false positives on startup.
package prescan

import (
	"fmt"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Result holds the outcome of a pre-scan pass.
type Result struct {
	Ports     []scanner.Port
	ScannedAt time.Time
	Duration  time.Duration
}

// Config controls pre-scan behaviour.
type Config struct {
	PortStart int
	PortEnd   int
	Retries   int
	Delay     time.Duration
}

// DefaultConfig returns sensible pre-scan defaults.
func DefaultConfig() Config {
	return Config{
		PortStart: 1,
		PortEnd:   65535,
		Retries:   3,
		Delay:     500 * time.Millisecond,
	}
}

// Run executes the pre-scan, retrying on error up to cfg.Retries times.
func Run(cfg Config) (Result, error) {
	var (
		ports []scanner.Port
		err   error
	)

	start := time.Now()

	for attempt := 0; attempt < cfg.Retries; attempt++ {
		ports, err = scanner.ScanPorts(cfg.PortStart, cfg.PortEnd)
		if err == nil {
			break
		}
		time.Sleep(cfg.Delay)
	}

	if err != nil {
		return Result{}, fmt.Errorf("prescan: all %d attempts failed: %w", cfg.Retries, err)
	}

	return Result{
		Ports:     ports,
		ScannedAt: time.Now(),
		Duration:  time.Since(start),
	}, nil
}
