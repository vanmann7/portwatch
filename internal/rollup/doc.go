// Package rollup provides event batching for portwatch.
//
// During rapid port-state churn (e.g. a service restart cycling dozens of
// ephemeral ports) individual per-port alerts create noise. Roller collects
// events within a configurable sliding window and emits a single Batch,
// allowing downstream handlers to produce one consolidated notification.
//
// Typical usage:
//
//	roller, batches := rollup.New(500 * time.Millisecond)
//	for batch := range batches {
//		fmt.Printf("batch: %d events\n", len(batch.Events))
//	}
package rollup
