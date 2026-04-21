// Package ratelimit provides a token-bucket rate limiter for controlling
// how frequently port-change notifications are emitted. Callers consume
// tokens on each event; tokens refill automatically after a configurable
// interval, preventing alert storms during rapid port churn.
package ratelimit
