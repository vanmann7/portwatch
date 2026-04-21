// Package watchdog provides a liveness monitor for long-running components.
//
// A Watchdog expects periodic Ping calls from the component it supervises.
// If pings stop arriving within the configured stale threshold the watchdog
// fires an onStale callback; if they continue to be absent past the dead
// threshold it fires an onDead callback.
//
// Typical usage:
//
//	w := watchdog.New(5*time.Second, 15*time.Second, onStale, onDead)
//	go w.Run(ctx, time.Second)
//	// inside the supervised component:
//	w.Ping()
package watchdog
