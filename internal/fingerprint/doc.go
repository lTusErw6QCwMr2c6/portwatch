// Package fingerprint provides stable hash-based identifiers for port alert
// events, enabling deduplication and correlation across cycles.
//
// A fingerprint is a short, deterministic string derived from the key
// attributes of a port event (e.g. protocol, host, port number, and state).
// Identical events observed in different cycles produce the same fingerprint,
// allowing callers to suppress duplicate alerts and track state changes over
// time without storing full event records.
package fingerprint
