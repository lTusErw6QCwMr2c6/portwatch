package eventcounter

import (
	"time"

	"github.com/user/portwatch/internal/alert"
)

// Handler wraps a Counter and fires a callback when a key exceeds a threshold
// within the rolling window.
type Handler struct {
	counter   *Counter
	threshold int
	onExceed  func(key string, count int)
}

// NewHandler creates a Handler that fires onExceed when the event count for a
// key surpasses threshold within window.
func NewHandler(window time.Duration, threshold int, onExceed func(key string, count int)) *Handler {
	return &Handler{
		counter:   New(window),
		threshold: threshold,
		onExceed:  onExceed,
	}
}

// Handle records the event and invokes the callback if the threshold is met.
func (h *Handler) Handle(e alert.Event) {
	key := e.Port.Protocol + ":" + itoa(e.Port.Number)
	h.counter.Record(key)
	count := h.counter.Count(key)
	if count >= h.threshold && h.onExceed != nil {
		h.onExceed(key, count)
	}
}

// Reset clears all counters in the underlying Counter.
func (h *Handler) Reset() {
	h.counter.Reset()
}

// itoa converts an int to its decimal string representation without importing
// strconv at the call site.
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	buf := [20]byte{}
	pos := len(buf)
	for n > 0 {
		pos--
		buf[pos] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[pos:])
}
