// Package sigterm provides graceful shutdown handling for portwatch.
package sigterm

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// Handler listens for OS termination signals and cancels a context.
type Handler struct {
	cancel context.CancelFunc
	ch     chan os.Signal
	done   chan struct{}
}

// New creates a Handler and starts listening for SIGINT and SIGTERM.
// The returned context is cancelled when a signal is received.
func New(parent context.Context) (context.Context, *Handler) {
	ctx, cancel := context.WithCancel(parent)
	h := &Handler{
		cancel: cancel,
		ch:     make(chan os.Signal, 1),
		done:   make(chan struct{}),
	}
	signal.Notify(h.ch, syscall.SIGINT, syscall.SIGTERM)
	go h.run()
	return ctx, h
}

func (h *Handler) run() {
	defer close(h.done)
	_, ok := <-h.ch
	if ok {
		h.cancel()
	}
}

// Stop unregisters signal handling and releases resources.
func (h *Handler) Stop() {
	signal.Stop(h.ch)
	close(h.ch)
	<-h.done
}
