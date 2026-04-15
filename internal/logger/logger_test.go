package logger

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestLog_FormatsCorrectly(t *testing.T) {
	var buf bytes.Buffer
	l := New(&buf)

	fixedTime := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	err := l.Log(Event{
		Timestamp: fixedTime,
		Level:     LevelInfo,
		Message:   "port 8080 opened",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := buf.String()
	want := "2024-06-01 12:00:00 [INFO] port 8080 opened\n"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestLog_Levels(t *testing.T) {
	tests := []struct {
		name    string
		level   Level
		contain string
	}{
		{"info level", LevelInfo, "[INFO]"},
		{"warn level", LevelWarn, "[WARN]"},
		{"error level", LevelError, "[ERROR]"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			l := New(&buf)
			_ = l.Log(Event{Timestamp: time.Now(), Level: tc.level, Message: "test"})
			if !strings.Contains(buf.String(), tc.contain) {
				t.Errorf("expected output to contain %q, got %q", tc.contain, buf.String())
			}
		})
	}
}

func TestNew_DefaultsToStdout(t *testing.T) {
	l := New(nil)
	if l.out == nil {
		t.Error("expected non-nil writer when nil is passed")
	}
}

func TestInfo_WritesMessage(t *testing.T) {
	var buf bytes.Buffer
	l := New(&buf)
	_ = l.Info("port 443 opened")
	if !strings.Contains(buf.String(), "port 443 opened") {
		t.Errorf("expected message in output, got %q", buf.String())
	}
	if !strings.Contains(buf.String(), "[INFO]") {
		t.Errorf("expected [INFO] level in output, got %q", buf.String())
	}
}
