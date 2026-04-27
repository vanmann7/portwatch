// Package portage tracks the continuous open duration of each monitored port
// and classifies ports into age buckets:
//
//   - new:          open for less than 5 minutes
//   - established:  open for 5–60 minutes
//   - long-running: open for more than 60 minutes
//
// Age classification helps operators quickly distinguish freshly-opened ports
// from long-lived services during incident triage.
package portage
