// Package metrics provides a thread-safe tracker for portwatch scan statistics.
//
// Use New to create a Tracker, then call RecordScan after each scan cycle.
// Snapshot returns a point-in-time copy of all accumulated counters without
// blocking ongoing scans.
package metrics
