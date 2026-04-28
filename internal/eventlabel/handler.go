package eventlabel

import (
	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/logger"
)

// Handler wraps a Labeler and invokes a downstream callback with the
// computed labels for each event.
type Handler struct {
	labeler  *Labeler
	log      *logger.Logger
	callback func(alert.Event, map[string]string)
}

// NewHandler creates a Handler that applies l to every event and passes
// the resulting label map to callback. callback is always called, even
// when the label map is nil.
func NewHandler(l *Labeler, log *logger.Logger, callback func(alert.Event, map[string]string)) *Handler {
	return &Handler{labeler: l, log: log, callback: callback}
}

// Handle applies the labeler rules to ev and forwards the event and
// labels to the registered callback.
func (h *Handler) Handle(ev alert.Event) {
	labels := h.labeler.Apply(ev)
	if len(labels) > 0 {
		h.log.Info("eventlabel: applied %d label(s) to event port=%d", len(labels), ev.Port)
	}
	h.callback(ev, labels)
}
