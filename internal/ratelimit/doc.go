// Package ratelimit provides a token-bucket rate limiter for controlling
// how frequently port-change alerts are emitted. Callers acquire a token
// before dispatching a notification; if the bucket is empty the event is
// dropped until tokens are refilled on the next interval tick.
package ratelimit
