package eventtag

import (
	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/logger"
)

// Handler wraps a Tagger and forwards tagged events to a callback.
type Handler struct {
	tagger *Tagger
	log    *logger.Logger
	next   func(ev alert.Event, tags []string)
}

// NewHandler creates a Handler that calls next for every event, passing
// the matched tags (which may be nil when no rules match).
func NewHandler(tagger *Tagger, log *logger.Logger, next func(alert.Event, []string)) *Handler {
	if next == nil {
		next = func(alert.Event, []string) {}
	}
	return &Handler{tagger: tagger, log: log, next: next}
}

// Handle processes a single event, resolves its tags, logs a summary when
// tags are found, and invokes the downstream callback.
func (h *Handler) Handle(ev alert.Event) {
	tags := h.tagger.Apply(ev)
	if len(tags) > 0 {
		h.log.Info("eventtag: port %d/%s matched %d tag(s)", ev.Port, ev.Protocol, len(tags))
	}
	h.next(ev, tags)
}
