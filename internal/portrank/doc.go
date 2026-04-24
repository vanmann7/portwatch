// Package portrank provides priority ranking for monitored ports.
//
// Ranks are integers in the range [1, 100] where higher values indicate
// greater priority. Built-in ranks cover common sensitive services; callers
// may supply overrides at construction time or at runtime via SetOverride.
//
// The Top helper selects the highest-ranked ports from an arbitrary list,
// useful for surfacing the most critical changes when many ports change at
// once.
package portrank
