package watcher

import (
	"time"

	"github.com/user/portwatch/internal/logger"
	"github.com/user/portwatch/internal/scanner"
)

// Config holds the watcher configuration.
type Config struct {
	Interval  time.Duration
	MinPort   int
	MaxPort   int
	Protocols []string
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Interval:  5 * time.Second,
		MinPort:   1,
		MaxPort:   65535,
		Protocols: []string{"tcp", "udp"},
	}
}

// Watcher monitors port activity and logs changes.
type Watcher struct {
	cfg    Config
	log    *logger.Logger
	stop   chan struct{}
}

// New creates a new Watcher with the given config and logger.
func New(cfg Config, log *logger.Logger) *Watcher {
	return &Watcher{
		cfg:  cfg,
		log:  log,
		stop: make(chan struct{}),
	}
}

// Start begins the polling loop. It blocks until Stop is called.
func (w *Watcher) Start() {
	w.log.Info("portwatch started", "interval", w.cfg.Interval)

	prev, err := scanner.ScanPorts(w.cfg.MinPort, w.cfg.MaxPort, w.cfg.Protocols)
	if err != nil {
		w.log.Error("initial scan failed", "err", err)
		return
	}

	ticker := time.NewTicker(w.cfg.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			curr, err := scanner.ScanPorts(w.cfg.MinPort, w.cfg.MaxPort, w.cfg.Protocols)
			if err != nil {
				w.log.Error("scan failed", "err", err)
				continue
			}

			opened, closed := scanner.Diff(prev, curr)

			for _, p := range opened {
				w.log.Info("port opened", "port", scanner.FormatPort(p))
			}
			for _, p := range closed {
				w.log.Info("port closed", "port", scanner.FormatPort(p))
			}

			prev = curr

		case <-w.stop:
			w.log.Info("portwatch stopped")
			return
		}
	}
}

// Stop signals the watcher to cease polling.
func (w *Watcher) Stop() {
	close(w.stop)
}
