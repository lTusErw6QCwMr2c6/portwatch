// Package limiter enforces a configurable cap on the number of simultaneously
// tracked open ports, dropping Opened events that would exceed the limit and
// decrementing the counter on Closed events.
package limiter
