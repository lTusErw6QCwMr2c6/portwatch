package config

import (
	"flag"
	"fmt"
	"os"
)

// Config holds the runtime configuration for portwatch.
type Config struct {
	Interval  int
	PortStart int
	PortEnd   int
	LogFile   string
	LogLevel  string
	JSON      bool
}

// DefaultConfig returns a Config populated with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		Interval:  5,
		PortStart: 1,
		PortEnd:   65535,
		LogFile:   "",
		LogLevel:  "info",
		JSON:      false,
	}
}

// Parse parses command-line flags and returns a populated Config.
func Parse() (*Config, error) {
	cfg := DefaultConfig()

	flag.IntVar(&cfg.Interval, "interval", cfg.Interval, "Scan interval in seconds")
	flag.IntVar(&cfg.PortStart, "port-start", cfg.PortStart, "Start of port range to monitor")
	flag.IntVar(&cfg.PortEnd, "port-end", cfg.PortEnd, "End of port range to monitor")
	flag.StringVar(&cfg.LogFile, "log-file", cfg.LogFile, "Path to log file (default: stdout)")
	flag.StringVar(&cfg.LogLevel, "log-level", cfg.LogLevel, "Log level: debug, info, warn, error")
	flag.BoolVar(&cfg.JSON, "json", cfg.JSON, "Output logs in JSON format")

	flag.Parse()

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate checks that the configuration values are acceptable.
func (c *Config) Validate() error {
	if c.Interval < 1 {
		return fmt.Errorf("interval must be at least 1 second, got %d", c.Interval)
	}
	if c.PortStart < 1 || c.PortStart > 65535 {
		return fmt.Errorf("port-start must be between 1 and 65535, got %d", c.PortStart)
	}
	if c.PortEnd < 1 || c.PortEnd > 65535 {
		return fmt.Errorf("port-end must be between 1 and 65535, got %d", c.PortEnd)
	}
	if c.PortStart > c.PortEnd {
		return fmt.Errorf("port-start (%d) must not exceed port-end (%d)", c.PortStart, c.PortEnd)
	}
	valid := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
	if !valid[c.LogLevel] {
		return fmt.Errorf("invalid log-level %q; must be one of: debug, info, warn, error", c.LogLevel)
	}
	if c.LogFile != "" {
		if _, err := os.Stat(c.LogFile); os.IsNotExist(err) {
			f, err := os.Create(c.LogFile)
			if err != nil {
				return fmt.Errorf("cannot create log file %q: %w", c.LogFile, err)
			}
			f.Close()
		}
	}
	return nil
}

// PortCount returns the total number of ports in the configured range.
func (c *Config) PortCount() int {
	return c.PortEnd - c.PortStart + 1
}
