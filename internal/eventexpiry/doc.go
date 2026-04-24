// Package eventexpiry implements a time-to-live registry for alert events.
// Events are tracked by an opaque string key; when the TTL elapses without
// the entry being cancelled the registered ExpiredFunc is invoked.
// Typical use-cases include detecting ports that opened but never closed
// within an expected window, or enforcing maximum observation durations.
package eventexpiry
