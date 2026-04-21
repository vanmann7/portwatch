// Package ratelimit provides a token-bucket rate limiter for controlling
// how frequently port change notifications are emitted. It prevents alert
// storms when many ports change state simultaneously.
package ratelimit
