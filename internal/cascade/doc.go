// Package cascade provides a composable multi-stage event processing
// chain for portwatch. Stages are applied in order and may filter,
// transform, or enrich alert events.
//
// A typical cascade pipeline looks like:
//
//	pipeline := cascade.New(
//		filter.ByPort(80, 443),
//		enrich.WithGeoIP(db),
//		alert.ToSlack(webhook),
//	)
//	pipeline.Run(ctx, events)
//
// Each stage receives the event produced by the previous stage. If a
// stage returns a nil event, processing stops and the event is dropped.
package cascade
