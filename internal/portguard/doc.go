// Package portguard evaluates port numbers against an ordered list of
// allow/deny/warn rules, returning a Verdict for each port.
//
// Rules are checked in declaration order; the first matching rule wins.
// If no rule matches, the Guard falls back to a configurable default Verdict.
//
// Rules can be loaded programmatically or parsed from a plain-text
// configuration file using ParseRules.
package portguard
