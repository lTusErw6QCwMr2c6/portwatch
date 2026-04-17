package export

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// FileHandler writes exported data to a file on each flush.
type FileHandler struct {
	dir    string
	format Format
}

// NewFileHandler returns a FileHandler that writes files to dir.
func NewFileHandler(dir string, f Format) *FileHandler {
	return &FileHandler{dir: dir, format: f}
}

// Flush exports events to a timestamped file in the configured directory.
func (h *FileHandler) Flush(events []alert.Event) error {
	if len(events) == 0 {
		return nil
	}
	if err := os.MkdirAll(h.dir, 0o755); err != nil {
		return fmt.Errorf("export: mkdir: %w", err)
	}
	ext := string(h.format)
	name := fmt.Sprintf("portwatch-%s.%s", time.Now().UTC().Format("20060102T150405Z"), ext)
	path := filepath.Join(h.dir, name)
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("export: create file: %w", err)
	}
	defer f.Close()
	ex := New(f, h.format)
	if err := ex.Write(events); err != nil {
		return fmt.Errorf("export: write: %w", err)
	}
	return nil
}
