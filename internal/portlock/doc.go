// Package portlock provides a thread-safe registry of operator-locked ports.
//
// Locked ports are intentionally pinned open (or closed) and should be
// excluded from unexpected-change alerts. Each lock entry carries an optional
// human-readable reason so operators can understand why a port was pinned.
//
// The lock table can be persisted to a JSON file and reloaded on startup so
// that locks survive daemon restarts.
package portlock
