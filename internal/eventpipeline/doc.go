// Package eventpipeline provides a composable, ordered processing pipeline for
// alert.Event slices. Each stage is a named transform function; stages are
// executed in registration order, with early exit when the event slice becomes
// empty.
package eventpipeline
