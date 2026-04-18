package anomaly

import (
	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/logger"
)

// Handler wraps a Detector and logs detections via a Logger.
type Handler struct {
	detector *Detector
	log      *logger.Logger
	onDetect func(Detection)
}

// NewHandler returns a Handler that logs anomalies and optionally calls onDetect.
func NewHandler(d *Detector, log *logger.Logger, onDetect func(Detection)) *Handler {
	if onDetect == nil {
		onDetect = func(Detection) {}
	}
	return &Handler{detector: d, log: log, onDetect: onDetect}
}

// Handle evaluates a slice of events and triggers callbacks for any anomalies found.
func (h *Handler) Handle(events []alert.Event) {
	for _, e := range events {
		if det, ok := h.detector.Evaluate(e); ok {
			h.log.Warn("anomaly detected: " + det.String())
			h.onDetect(det)
		}
	}
}
