// Package healthcheck provides a lightweight HTTP server that exposes
// a /healthz endpoint for the portwatch daemon.
//
// The endpoint returns a JSON payload containing the current daemon
// status, last scan timestamp, and accumulated metrics so that external
// monitoring tools (e.g. uptime checkers, Kubernetes liveness probes)
// can verify the daemon is running and scanning as expected.
package healthcheck
