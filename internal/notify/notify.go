package notify

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// Channel represents a notification destination.
type Channel interface {
	Send(event alert.Event) error
}

// StdoutChannel writes notifications to an io.Writer.
type StdoutChannel struct {
	Writer io.Writer
	Prefix string
}

// NewStdoutChannel returns a StdoutChannel writing to os.Stdout.
func NewStdoutChannel(prefix string) *StdoutChannel {
	return &StdoutChannel{Writer: os.Stdout, Prefix: prefix}
}

// Send writes a formatted event line to the writer.
func (c *StdoutChannel) Send(event alert.Event) error {
	ts := time.Now().Format(time.RFC3339)
	_, err := fmt.Fprintf(c.Writer, "%s [%s] %s\n", ts, c.Prefix, event.String())
	return err
}

// Dispatcher fans out events to one or more channels.
type Dispatcher struct {
	channels []Channel
}

// NewDispatcher creates a Dispatcher with the given channels.
func NewDispatcher(channels ...Channel) *Dispatcher {
	return &Dispatcher{channels: channels}
}

// Dispatch sends the event to all registered channels, collecting errors.
func (d *Dispatcher) Dispatch(event alert.Event) []error {
	var errs []error
	for _, ch := range d.channels {
		if err := ch.Send(event); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

// DispatchAll sends multiple events through the dispatcher.
func (d *Dispatcher) DispatchAll(events []alert.Event) []error {
	var errs []error
	for _, e := range events {
		errs = append(errs, d.Dispatch(e)...)
	}
	return errs
}
