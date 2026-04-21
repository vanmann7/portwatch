// Package ratelimit provides a token-bucket rate limiter for controlling
// how frequently port-change alerts are emitted. Tokens are consumed on
// each allowed event and refilled automatically after a configurable interval.
package ratelimit
