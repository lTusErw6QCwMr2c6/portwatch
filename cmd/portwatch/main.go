package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/logger"
	"github.com/user/portwatch/internal/watcher"
)

func main() {
	cfg, err := config.Parse()
	if err != nil {
		fmt.Fprintf(os.Stderr, "portwatch: configuration error: %v\n", err)
		os.Exit(1)
	}

	var logDest *os.File
	if cfg.LogFile != "" {
		f, err := os.OpenFile(cfg.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "portwatch: failed to open log file: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()
		logDest = f
	} else {
		logDest = os.Stdout
	}

	log := logger.New(logDest, cfg.LogLevel)
	log.Info("portwatch starting", "interval", cfg.Interval, "ports", fmt.Sprintf("%d-%d", cfg.PortStart, cfg.PortEnd))

	watcherCfg := watcher.DefaultConfig()
	watcherCfg.Interval = cfg.Interval
	watcherCfg.PortStart = cfg.PortStart
	watcherCfg.PortEnd = cfg.PortEnd

	w, err := watcher.New(watcherCfg, log)
	if err != nil {
		log.Error("failed to create watcher", "error", err)
		os.Exit(1)
	}

	w.Start()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	log.Info("portwatch shutting down")
	w.Stop()
}
