package eventtransform

import (
	"fmt"

	"github.com/user/portwatch/internal/alert"
)

// Handler wraps a Transformer and a downstream callback, applying all
// registered stages before forwarding the result.
type Handler struct {
	t    *Transformer
	next func(alert.Event) error
}

// NewHandler creates a Handler that transforms events using t before passing
// them to next. Both arguments must be non-nil.
func NewHandler(t *Transformer, next func(alert.Event) error) (*Handler, error) {
	if t == nil {
		return nil, fmt.Errorf("eventtransform: transformer must not be nil")
	}
	if next == nil {
		return nil, fmt.Errorf("eventtransform: next handler must not be nil")
	}
	return &Handler{t: t, next: next}, nil
}

// Handle transforms e and forwards the result to the downstream handler.
func (h *Handler) Handle(e alert.Event) error {
	out, err := h.t.Apply(e)
	if err != nil {
		return err
	}
	return h.next(*out)
}
