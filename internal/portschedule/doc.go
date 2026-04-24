// Package portschedule enforces per-port minimum scan intervals so that
// high-frequency watch loops do not redundantly re-scan ports that were
// checked recently.
//
// Usage:
//
//	sched := portschedule.New(30 * time.Second)
//	if sched.Due(port) {
//	    // perform scan
//	    sched.Record(port)
//	}
package portschedule
