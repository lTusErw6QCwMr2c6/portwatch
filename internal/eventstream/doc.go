// Package eventstream implements a fan-out broadcast stream for alert.Event
// values. Multiple named consumers can subscribe and receive events
// independently; slow consumers experience drop-on-full backpressure rather
// than blocking the publisher.
package eventstream
