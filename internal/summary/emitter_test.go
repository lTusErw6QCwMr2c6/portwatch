package summary_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/logger"
	"github.com/user/portwatch/internal/summary"
)

func TestEmitter_LogsSummary(t *testing.T) {
	var buf bytes.Buffer
	log := logger.New(&buf, logger.Info)
	m := makeMetrics(1, 2, 0, 0)
	b := summary.New(m, 5)
	e := summary.NewEmitter(b, 30*time.Millisecond, log)
	e.Start(func() []alert.Event { return nil })
	defer e.Stop()

	time.Sleep(60 * time.Millisecond)

	if !strings.Contains(buf.String(), "Summary") {
		t.Errorf("expected summary in log output, got: %s", buf.String())
	}
}

func TestEmitter_Stop_Halts(t *testing.T) {
	var buf bytes.Buffer
	log := logger.New(&buf, logger.Info)
	m := makeMetrics(0, 0, 0, 0)
	b := summary.New(m, 5)
	e := summary.NewEmitter(b, 20*time.Millisecond, log)
	e.Start(func() []alert.Event { return nil })
	e.Stop()

	time.Sleep(50 * time.Millisecond)
	first := buf.Len()
	time.Sleep(30 * time.Millisecond)
	if buf.Len() != first {
		t.Error("emitter continued after Stop")
	}
}
