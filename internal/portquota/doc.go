// Package portquota provides a rolling-window quota tracker for per-port alert
// rate limiting.
//
// A Tracker is created with a maximum alert count and a window duration. Each
// call to Allow increments the counter for the given port. Once the counter
// reaches the maximum, Allow returns false until the window expires and the
// counter resets automatically.
//
// Typical usage:
//
//	q := portquota.New(5, time.Minute)
//	if q.Allow(port) {
//	    // emit alert
//	}
package portquota
