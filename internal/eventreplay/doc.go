// Package eventreplay provides a bounded in-memory buffer that records
// alert.Event values so they can be replayed or inspected after the fact.
// It is safe for concurrent use and evicts the oldest entry when the
// configured capacity is exceeded.
package eventreplay
