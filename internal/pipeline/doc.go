// Package pipeline wires together the filter, rate-limit, and alert stages
// of the portwatch event processing chain.
//
// A Pipeline accepts a scanner.Diff, optionally filters ports, checks alert
// thresholds, and suppresses duplicate events via rate-limiting before
// returning the final []alert.Event slice to the caller.
package pipeline
