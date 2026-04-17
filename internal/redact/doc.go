// Package redact scrubs sensitive data from strings before they are
// written to logs, audit files, or exported snapshots. Callers register
// one or more Rules (each a compiled regexp and a replacement string)
// and call Apply / ApplyAll on any text that may contain credentials,
// IP addresses, or other private information.
package redact
