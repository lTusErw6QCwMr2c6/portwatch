// Package audit provides a persistent audit trail of port activity events.
package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// Entry represents a single audited event written to the audit log.
type Entry struct {
	Timestamp time.Time   `json:"timestamp"`
	Event     alert.Event `json:"event"`
}

// Writer appends audit file in newline-delimited JSON format.
type Writer struct {
	mu   sync.Mutex
	file *os.File
}

// New opens or creates the audit log file at the given path.
func New(path{
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, fmt.Errorf("audit: open %s: %w", path, err)
	}
	return &Writer{file: f}, nil
}

// Write encodes the event as a JSON entry and appends it to the audit log.
func (w *Writer) Write(e alert.Event) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	entry := Entry{Timestamp: time.Now().UTC(), Event: e}
	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("audit: marshal: %w", err)
	}
	data = append(data, '\n')
	_, err = w.file.Write(data)
	return err
}

// Close closes the underlying file.
func (w *Writer) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.file.Close()
}
