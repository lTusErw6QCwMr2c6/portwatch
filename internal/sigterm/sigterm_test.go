package sigterm_test

import (
	"context"
	"syscall"
	"testing"
	"time"

	"github.com/user/portwatch/internal/sigterm"
)

func TestNew_ReturnsNonNilContext(t *testing.T) {
	ctx, h := sigterm.New(context.Background())
	defer h.Stop()
	if ctx == nil {
		t.Fatal("expected non-nil context")
	}
}

func TestStop_DoesNotBlock(t *testing.T) {
	_, h := sigterm.New(context.Background())
	done := make(chan struct{})
	go func() {
		h.Stop()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("Stop blocked")
	}
}

func TestSignal_CancelsContext(t *testing.T) {
	ctx, h := sigterm.New(context.Background())
	defer h.Stop()

	p, err := os.FindProcess(os.Getpid())
	if err != nil {
		t.Fatalf("find process: %v", err)
	}
	if err := p.Signal(syscall.SIGINT); err != nil {
		t.Fatalf("send signal: %v", err)
	}

	select {
	case <-ctx.Done():
	case <-time.After(time.Second):
		t.Fatal("context not cancelled after SIGINT")
	}
}
