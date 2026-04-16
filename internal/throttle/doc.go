// Package throttle provides a sliding-window rate limiter keyed by string
// identifiers. Unlike ratelimit (which enforces a per-event cooldown),
// throttle tracks a maximum number of occurrences within a rolling time
// window, making it suitable for burst-aware suppression of repeated port
// events before they reach the alert or notify pipeline.
package throttle
