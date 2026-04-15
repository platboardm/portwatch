// Package dedup implements alert deduplication for portwatch.
//
// A Deduplicator tracks a fingerprint of each alert it processes and suppresses
// repeated identical alerts that arrive within a configurable time window. This
// prevents alert storms when a service is flapping rapidly.
//
// Fingerprints are derived from the alert's target name, host, port and
// severity using a SHA-256 hash, so only semantically identical alerts are
// collapsed — distinct services or severity changes are always forwarded.
//
// Usage:
//
//	d := dedup.New(5 * time.Minute)
//	if d.Allow(a) {
//		// forward alert
//	}
//
// Deduplicator is safe for concurrent use.
package dedup
