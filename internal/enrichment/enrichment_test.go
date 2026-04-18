package enrichment_test

import (
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/enrichment"
	"github.com/user/portwatch/internal/process"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/tag"
)

func makeEvent(port int, proto, state string) alert.Event {
	return alert.Event{
		Port:  scanner.Port{Port: port, Proto: proto},
		State: state,
	}
}

func noopLookup(_ int) process.Info {
	return process.Info{PID: 0, Name: ""}
}

func TestEnrich_AttachesTags(t *testing.T) {
	tagger := tag.New(nil)
	e := enrichment.New(tagger, noopLookup)
	ev := makeEvent(80, "tcp", "opened")
	m := e.Enrich(ev)
	if len(m.Tags) == 0 {
		t.Error("expected tags for port 80")
	}
	if m.Note == "" {
		t.Error("expected non-empty note for well-known port")
	}
}

func TestEnrich_UnknownPort_NoTags(t *testing.T) {
	tagger := tag.New(nil)
	e := enrichment.New(tagger, noopLookup)
	ev := makeEvent(39999, "tcp", "opened")
	m := e.Enrich(ev)
	if len(m.Tags) != 0 {
		t.Errorf("expected no tags, got %v", m.Tags)
	}
	if m.Note != "" {
		t.Errorf("expected empty note, got %q", m.Note)
	}
}

func TestEnrich_AttachesProcess(t *testing.T) {
	tagger := tag.New(nil)
	lookup := func(_ int) process.Info { return process.Info{PID: 42, Name: "nginx"} }
	e := enrichment.New(tagger, lookup)
	m := e.Enrich(makeEvent(8080, "tcp", "opened"))
	if m.Process == "" {
		t.Error("expected non-empty process string")
	}
}

func TestEnrichAll_ReturnsOnePerEvent(t *testing.T) {
	tagger := tag.New(nil)
	e := enrichment.New(tagger, noopLookup)
	events := []alert.Event{
		makeEvent(80, "tcp", "opened"),
		makeEvent(443, "tcp", "opened"),
		makeEvent(22, "tcp", "closed"),
	}
	results := e.EnrichAll(events)
	if len(results) != len(events) {
		t.Errorf("expected %d results, got %d", len(events), len(results))
	}
}
