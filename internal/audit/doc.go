// Package audit records a persistent, append-only log of port activity
// events detected by portwatch. Each entry is written as newline-delimited
// JSON and includes a UTC timestamp alongside the alert event payload.
package audit
