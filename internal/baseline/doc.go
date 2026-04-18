// Package baseline provides a persistent set of "known-good" ports.
//
// The Baseline type is safe for concurrent use. Callers load a previously
// saved baseline at startup, compare live scan results against it, and
// optionally persist an updated baseline after operator review.
package baseline
