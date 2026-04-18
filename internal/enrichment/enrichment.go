// Package enrichment attaches contextual metadata to alert events.
package enrichment

import (
	"fmt"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/tag"
	"github.com/user/portwatch/internal/process"
)

// Metadata holds enriched context for an event.
type Metadata struct {
	Event   alert.Event
	Tags    []string
	Process string
	Note    string
}

// Enricher attaches tags and process info to events.
type Enricher struct {
	tagger  *tag.Tagger
	lookup  func(port int) process.Info
}

// New returns an Enricher using the provided tagger and process lookup func.
func New(t *tag.Tagger, lookup func(port int) process.Info) *Enricher {
	return &Enricher{tagger: t, lookup: lookup}
}

// Enrich returns a Metadata record for the given event.
func (e *Enricher) Enrich(ev alert.Event) Metadata {
	tags := e.tagger.LookupAll(ev.Port.Port, ev.Port.Proto)
	info := e.lookup(ev.Port.Port)
	note := ""
	if len(tags) > 0 {
		note = fmt.Sprintf("well-known: %s", tags[0])
	}
	return Metadata{
		Event:   ev,
		Tags:    tags,
		Process: info.String(),
		Note:    note,
	}
}

// EnrichAll enriches a slice of events.
func (e *Enricher) EnrichAll(events []alert.Event) []Metadata {
	out := make([]Metadata, 0, len(events))
	for _, ev := range events {
		out = append(out, e.Enrich(ev))
	}
	return out
}
