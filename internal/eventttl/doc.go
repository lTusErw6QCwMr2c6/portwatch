// Package eventttl implements a time-to-live store for port events.
// It allows callers to track recently seen events and automatically
// evict stale ones once their configured lifetime has elapsed.
package eventttl
