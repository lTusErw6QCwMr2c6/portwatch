// Package eventprojector provides a latest-event projection over a
// stream of alert.Event values, keyed by protocol+address. It is
// useful for building a live view of currently open ports without
// replaying the full event history.
package eventprojector
