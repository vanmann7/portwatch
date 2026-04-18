// Package throttle provides a sliding-window call throttle for use in
// portwatch to limit the frequency of alerts and notifications.
//
// A Throttle allows at most N calls within a given time window, making it
// useful for suppressing alert floods when many ports change at once.
package throttle
