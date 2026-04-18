package anomaly

import (
	"bytes"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/logger"
)

func silentLogger() *logger.Logger {
	return logger.New(&bytes.Buffer{}, logger.LevelInfo)
}

func TestNewHandler_NotNil(t *testing.T) {
	d := New(time.Second, 3)
	h := NewHandler(d, silentLogger(), nil)
	if h == nil {
		t.Fatal("expected non-nil handler")
	}
}

func TestHandle_NoAnomaly_CallbackNotFired(t *testing.T) {
	d := New(time.Second, 10)
	called := false
	h := NewHandler(d, silentLogger(), func(Detection) { called = true })
	h.Handle([]alert.Event{makeEvent("tcp", 80)})
	if called {
		t.Error("callback should not fire below threshold")
	}
}

func TestHandle_Anomaly_CallbackFired(t *testing.T) {
	d := New(time.Second, 2)
	called := 0
	h := NewHandler(d, silentLogger(), func(Detection) { called++ })
	events := []alert.Event{
		makeEvent("tcp", 1234),
		makeEvent("tcp", 1234),
		makeEvent("tcp", 1234),
	}
	h.Handle(events)
	if called == 0 {
		t.Error("expected callback to fire")
	}
}
