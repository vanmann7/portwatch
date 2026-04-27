// Package portbudget enforces configurable limits on the number of
// simultaneously open ports within named ranges.
//
// A Budget is created with one or more ranges, each specifying a low port,
// high port, and maximum allowed concurrent open ports. Callers update the
// budget via Open and Close as port state changes are detected, then call
// Exceeded to determine whether any range has breached its limit.
//
// Example:
//
//	b, _ := portbudget.New(
//		portbudget.WithRange(8000, 8999, 10),
//	)
//	b.Open(8080)
//	if ok, desc := b.Exceeded(); ok {
//		log.Printf("budget exceeded: %s", desc)
//	}
package portbudget
