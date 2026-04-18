package watchdog_test

import (
	"sync"
	"testing"
	"time"

	"github.com/user/portwatch/internal/watchdog"
)

type testLogger struct {
	mu   sync.Mutex
	msgs []string
}

func (l *testLogger) Warn(msg string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.msgs = append(l.msgs, msg)
}

func (l *testLogger) count() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return len(l.msgs)
}

func TestNew_NotStalled(t *testing.T) {
	log := &testLogger{}
	w := watchdog.New(500*time.Millisecond, log)
	defer w.Stop()
	if w.Stalled() {
		t.Fatal("expected not stalled immediately after creation")
	}
}

func TestStall_FiresAfterTimeout(t *testing.T) {
	log := &testLogger{}
	w := watchdog.New(80*time.Millisecond, log)
	defer w.Stop()
	time.Sleep(150 * time.Millisecond)
	if !w.Stalled() {
		t.Fatal("expected watchdog to be stalled after timeout")
	}
	if log.count() == 0 {
		t.Fatal("expected warning to be logged")
	}
}

func TestHeartbeat_ResetsStall(t *testing.T) {
	log := &testLogger{}
	w := watchdog.New(100*time.Millisecond, log)
	defer w.Stop()
	time.Sleep(60 * time.Millisecond)
	w.Heartbeat()
	time.Sleep(60 * time.Millisecond)
	if w.Stalled() {
		t.Fatal("expected heartbeat to prevent stall")
	}
}

func TestStop_PreventsStall(t *testing.T) {
	log := &testLogger{}
	w := watchdog.New(80*time.Millisecond, log)
	w.Stop()
	time.Sleep(150 * time.Millisecond)
	if w.Stalled() {
		t.Fatal("expected stopped watchdog to never stall")
	}
}
