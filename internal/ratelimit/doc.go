// Package ratelimit provides a token-bucket rate limiter for controlling
// the frequency of port-change notifications. A Limiter is initialised with
// a maximum number of tokens and a refill interval; each call to Allow
// consumes one token, and tokens are replenished automatically after the
// configured interval elapses.
package ratelimit
