package logger

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Level represents the severity of a log message.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelError Level = "ERROR"
)

// Event represents a port activity change event to be logged.
type Event struct {
	Timestamp time.Time
	Level     Level
	Message   string
}

// Logger writes structured port activity events to an output sink.
type Logger struct {
	out io.Writer
	timeFormat string
}

// New creates a Logger that writes to the given writer.
// If w is nil, os.Stdout is used.
func New(w io.Writer) *Logger {
	if w == nil {
		w = os.Stdout
	}
	return &Logger{
		out:        w,
		timeFormat: "2006-01-02 15:04:05",
	}
}

// Log writes a single Event to the output sink.
func (l *Logger) Log(e Event) error {
	line := fmt.Sprintf("%s [%s] %s\n",
		e.Timestamp.Format(l.timeFormat),
		e.Level,
		e.Message,
	)
	_, err := fmt.Fprint(l.out, line)
	return err
}

// Info logs an informational message with the current timestamp.
func (l *Logger) Info(msg string) error {
	return l.Log(Event{Timestamp: time.Now(), Level: LevelInfo, Message: msg})
}

// Warn logs a warning message with the current timestamp.
func (l *Logger) Warn(msg string) error {
	return l.Log(Event{Timestamp: time.Now(), Level: LevelWarn, Message: msg})
}

// Error logs an error message with the current timestamp.
func (l *Logger) Error(msg string) error {
	return l.Log(Event{Timestamp: time.Now(), Level: LevelError, Message: msg})
}
