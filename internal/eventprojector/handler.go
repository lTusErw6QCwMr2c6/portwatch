package eventprojector

import (
	"github.com/user/portwatch/internal/alert"
)

// Handler wraps a Projector and exposes a Handle method compatible
// with dispatch or pipeline handler signatures.
type Handler struct {
	proj     *Projector
	onUpdate func(alert.Event)
}

// NewHandler returns a Handler backed by the given Projector.
// onUpdate is an optional callback invoked after each projection
// update; pass nil to disable.
func NewHandler(proj *Projector, onUpdate func(alert.Event)) *Handler {
	if proj == nil {
		proj = New()
	}
	return &Handler{proj: proj, onUpdate: onUpdate}
}

// Handle applies the event to the projection and fires the optional
// callback.
func (h *Handler) Handle(e alert.Event) {
	h.proj.Apply(e)
	if h.onUpdate != nil {
		h.onUpdate(e)
	}
}

// Projector returns the underlying Projector for read access.
func (h *Handler) Projector() *Projector {
	return h.proj
}
