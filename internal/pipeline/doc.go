// Package pipeline wires the core portwatch processing stages together.
//
// It connects filter, dedupe, alert, and notify into an ordered chain so
// that callers only need to hand a state.Diff to Pipeline.Process and the
// rest is handled automatically.
package pipeline
