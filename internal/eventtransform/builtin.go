package eventtransform

import (
	"strings"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// UppercaseProtocol returns a TransformFunc that uppercases the event protocol.
func UppercaseProtocol() TransformFunc {
	return func(e *alert.Event) (*alert.Event, error) {
		e.Protocol = strings.ToUpper(e.Protocol)
		return e, nil
	}
}

// StampNow returns a TransformFunc that sets the event timestamp to the
// current UTC time if the event timestamp is zero.
func StampNow() TransformFunc {
	return func(e *alert.Event) (*alert.Event, error) {
		if e.Timestamp.IsZero() {
			e.Timestamp = time.Now().UTC()
		}
		return e, nil
	}
}

// SetLabel returns a TransformFunc that sets a fixed label key/value on the
// event's Labels map, allocating the map if necessary.
func SetLabel(key, value string) TransformFunc {
	return func(e *alert.Event) (*alert.Event, error) {
		if e.Labels == nil {
			e.Labels = make(map[string]string)
		}
		e.Labels[key] = value
		return e, nil
	}
}
