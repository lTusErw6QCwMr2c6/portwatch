package dispatch

import (
	"log"

	"github.com/user/portwatch/internal/alert"
)

// Middleware wraps a Handler with additional behaviour.
type Middleware func(Handler) Handler

// WithLogging returns a Middleware that logs each dispatched event.
func WithLogging(l *log.Logger) Middleware {
	return func(next Handler) Handler {
		return func(event alert.Event) error {
			l.Printf("dispatch: event type=%s port=%d/%s",
				event.Type, event.Port.Number, event.Port.Protocol)
			return next(event)
		}
	}
}

// WithRecovery returns a Middleware that recovers from panics in handlers.
func WithRecovery(l *log.Logger) Middleware {
	return func(next Handler) Handler {
		return func(event alert.Event) (err error) {
			defer func() {
				if r := recover(); r != nil {
					l.Printf("dispatch: recovered panic: %v", r)
				}
			}()
			return next(event)
		}
	}
}

// Apply wraps a Handler with a chain of middleware (first applied outermost).
func Apply(h Handler, mw ...Middleware) Handler {
	for i := len(mw) - 1; i >= 0; i-- {
		h = mw[i](h)
	}
	return h
}
