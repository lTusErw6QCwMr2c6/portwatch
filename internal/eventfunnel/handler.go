package eventfunnel

import (
	"github.com/example/portwatch/internal/alert"
)

// Handler is a function that processes a merged event.
type Handler func(ev alert.Event)

// Drain reads from the funnel's output channel and calls h for each event
// until the channel is closed. It blocks the calling goroutine.
func Drain(f *Funnel, h Handler) {
	for ev := range f.Out() {
		h(ev)
	}
}

// DrainAsync starts a goroutine that drains the funnel and calls h for each
// event. It returns a done channel that is closed when draining finishes.
func DrainAsync(f *Funnel, h Handler) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)
		Drain(f, h)
	}()
	return done
}
