// Package anomaly provides burst-based anomaly detection for port events.
// It tracks event frequency per port within a sliding time window and
// emits a Detection when activity exceeds a configured threshold.
package anomaly
