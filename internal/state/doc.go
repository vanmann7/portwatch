// Package state provides snapshot persistence and diff computation for
// portwatch. A Snapshot captures the set of open ports observed at a
// specific moment. Compare produces a Diff describing which ports were
// opened or closed between two consecutive snapshots, enabling the
// daemon to alert on unexpected changes.
package state
