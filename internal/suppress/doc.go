// Package suppress manages a suppression list of ports whose open/close
// events should be silently ignored. This is useful for well-known services
// that are always expected to be running and would otherwise generate
// constant noise in alerts.
package suppress
