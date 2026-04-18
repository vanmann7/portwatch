// Package ratelimit implements a token-bucket rate limiter used by portwatch
// to suppress alert floods when a large number of port changes are detected
// within a short time window.
//
// Typical usage:
//
//	limiter := ratelimit.New(10, time.Minute)
//	if limiter.Allow() {
//		// send alert
//	}
package ratelimit
