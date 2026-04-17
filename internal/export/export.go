// Package export writes port activity snapshots to external formats.
package export

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// Format represents an output format for exported data.
type Format string

const (
	FormatJSON Format = "json"
	FormatCSV  Format = "csv"
)

// Record is a flat representation of an alert event for export.
type Record struct {
	Timestamp time.Time `json:"timestamp"`
	Type      string    `json:"type"`
	Port      int       `json:"port"`
	Protocol  string    `json:"protocol"`
	Process   string    `json:"process"`
}

// Exporter writes records to an io.Writer in a given format.
type Exporter struct {
	format Format
	w      io.Writer
}

// New returns an Exporter that writes to w using the given format.
func New(w io.Writer, f Format) *Exporter {
	return &Exporter{format: f, w: w}
}

// Write encodes events into the configured format.
func (e *Exporter) Write(events []alert.Event) error {
	records := make([]Record, len(events))
	for i, ev := range events {
		records[i] = Record{
			Timestamp: ev.Timestamp,
			Type:      ev.Type.String(),
			Port:      ev.Port.Port,
			Protocol:  ev.Port.Protocol,
			Process:   ev.Port.Process.Name,
		}
	}
	switch e.format {
	case FormatJSON:
		return e.writeJSON(records)
	case FormatCSV:
		return e.writeCSV(records)
	default:
		return fmt.Errorf("unknown format: %s", e.format)
	}
}

func (e *Exporter) writeJSON(records []Record) error {
	enc := json.NewEncoder(e.w)
	enc.SetIndent("", "  ")
	return enc.Encode(records)
}

func (e *Exporter) writeCSV(records []Record) error {
	w := csv.NewWriter(e.w)
	if err := w.Write([]string{"timestamp", "type", "port", "protocol", "process"}); err != nil {
		return err
	}
	for _, r := range records {
		row := []string{
			r.Timestamp.Format(time.RFC3339),
			r.Type,
			fmt.Sprintf("%d", r.Port),
			r.Protocol,
			r.Process,
		}
		if err := w.Write(row); err != nil {
			return err
		}
	}
	w.Flush()
	return w.Error()
}
