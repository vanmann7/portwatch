package pipeline_test

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/dedupe"
	"github.com/user/portwatch/internal/filter"
	"github.com/user/portwatch/internal/metrics"
	"github.com/user/portwatch/internal/notify"
	"github.com/user/portwatch/internal/pipeline"
	"github.com/user/portwatch/internal/state"
)

func buildPipeline(t *testing.T, buf *bytes.Buffer, rules []string) *pipeline.Pipeline {
	t.Helper()
	f, err := filter.New(rules)
	if err != nil {
		t.Fatalf("filter.New: %v", err)
	}
	d := dedupe.New(5 * time.Minute)
	a := alert.New(buf)
	ch := notify.NewStdout(buf)
	disp := notify.NewDispatcher([]notify.Channel{ch})
	m := metrics.New()
	return pipeline.New(f, d, a, disp, m)
}

func TestProcessOpenedPort(t *testing.T) {
	var buf bytes.Buffer
	p := buildPipeline(t, &buf, nil)
	diff := state.Diff{Opened: []int{8080}, Closed: nil}
	if err := p.Process(context.Background(), diff); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() == 0 {
		t.Error("expected output in buffer, got none")
	}
}

func TestProcessFilteredPort(t *testing.T) {
	var buf bytes.Buffer
	// allow only port 9090; 8080 should be filtered out
	p := buildPipeline(t, &buf, []string{"9090"})
	diff := state.Diff{Opened: []int{8080}, Closed: nil}
	if err := p.Process(context.Background(), diff); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected no output for filtered port, got %q", buf.String())
	}
}

func TestProcessDeduplicatesEvent(t *testing.T) {
	var buf bytes.Buffer
	p := buildPipeline(t, &buf, nil)
	diff := state.Diff{Opened: []int{8080}}
	_ = p.Process(context.Background(), diff)
	buf.Reset()
	// second call with same event should be deduplicated
	if err := p.Process(context.Background(), diff); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected deduplicated event to produce no output, got %q", buf.String())
	}
}

func TestProcessNoDiff(t *testing.T) {
	var buf bytes.Buffer
	p := buildPipeline(t, &buf, nil)
	diff := state.Diff{}
	if err := p.Process(context.Background(), diff); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected no output for empty diff, got %q", buf.String())
	}
}
