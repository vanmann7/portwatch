// Package portcache implements a thread-safe, TTL-based in-memory cache for
// port scan results.
//
// Use [New] to create a cache with a desired TTL. Call [Cache.Set] after each
// probe and [Cache.Get] before probing to skip redundant network round-trips.
// Call [Cache.Evict] periodically (e.g. from the watch loop) to reclaim memory
// from stale entries.
package portcache
