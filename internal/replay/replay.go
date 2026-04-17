// Package replay provides functionality to replay historical port events
// from the audit log for debugging and analysis purposes.
package replay

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// Entry represents a single decoded audit log line.
type Entry struct {
	Timestamp time.Time  `json:"timestamp"`
	Event     alert.Event `json:"event"`
}

// Options controls replay behaviour.
type Options struct {
	Since  time.Time
	Until  time.Time
	Limit  int
}

// Reader reads and filters entries from an audit log file.
type Reader struct {
	path string
}

// New returns a Reader for the given audit log path.
func New(path string) *Reader {
	return &Reader{path: path}
}

// Read opens the audit file and returns entries matching the given options.
func (r *Reader) Read(opts Options) ([]Entry, error) {
	f, err := os.Open(r.path)
	if err != nil {
		return nil, fmt.Errorf("replay: open %s: %w", r.path, err)
	}
	defer f.Close()
	return readEntries(f, opts)
}

func readEntries(r io.Reader, opts Options) ([]Entry, error) {
	var results []Entry
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var e Entry
		if err := json.Unmarshal(line, &e); err != nil {
			continue
		}
		if !opts.Since.IsZero() && e.Timestamp.Before(opts.Since) {
			continue
		}
		if !opts.Until.IsZero() && e.Timestamp.After(opts.Until) {
			continue
		}
		results = append(results, e)
		if opts.Limit > 0 && len(results) >= opts.Limit {
			break
		}
	}
	return results, scanner.Err()
}
