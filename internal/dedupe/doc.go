// Package dedupe provides a Deduper that suppresses repeated port-change
// notifications for the same (port, event) pair within a sliding time window.
// This prevents alert storms when a port flaps rapidly.
package dedupe
