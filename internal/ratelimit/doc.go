// Package ratelimit provides a token-bucket rate limiter for controlling
// how frequently port-change alerts are emitted. Callers consume tokens
// via Allow and tokens are refilled automatically after the configured
// interval, preventing alert storms during rapid port-state changes.
package ratelimit
