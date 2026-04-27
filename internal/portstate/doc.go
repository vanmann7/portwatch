// Package portstate tracks per-port open/closed state and produces
// Transition values whenever a port changes state.
//
// Usage:
//
//	tr := portstate.New()
//	if trans, changed := tr.Update(80, portstate.Open); changed {
//		fmt.Printf("port %d: %v -> %v\n", trans.Port, trans.Prev, trans.Next)
//	}
package portstate
