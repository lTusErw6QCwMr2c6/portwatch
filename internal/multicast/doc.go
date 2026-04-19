// Package multicast implements a fan-out event bus that delivers port alert
// events to multiple named subscribers concurrently. Each subscriber receives
// its own buffered channel; slow consumers are skipped rather than blocking
// the publisher.
package multicast
