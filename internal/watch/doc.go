// Package watch implements the periodic monitoring loop for portwatch.
//
// A Watcher is constructed from a [config.Config] and an [alert.Alerter].
// Calling Run blocks until the provided context is cancelled, scanning
// the configured port range at each tick and emitting alerts whenever
// the set of open ports changes.
package watch
