// Package portpolicy provides an ordered, named allow/deny rule engine for
// port numbers. Rules are evaluated in insertion order; the first matching
// rule wins. If no rule matches, a configurable default action is applied.
//
// Rules can be loaded programmatically via AddRule or parsed from a plain-text
// configuration using ParseRules. Each rule covers a single port or a
// contiguous port range and carries a human-readable name that is included in
// the returned Verdict for observability.
package portpolicy
