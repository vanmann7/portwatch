package labeler_test

import (
	"testing"

	"github.com/yourorg/portwatch/internal/labeler"
)

func TestBuiltinSSH(t *testing.T) {
	l := labeler.New(nil)
	if got := l.Label(22); got != "ssh" {
		t.Fatalf("expected ssh, got %q", got)
	}
}

func TestBuiltinHTTPS(t *testing.T) {
	l := labeler.New(nil)
	if got := l.Label(443); got != "https" {
		t.Fatalf("expected https, got %q", got)
	}
}

func TestUnknownPortFallback(t *testing.T) {
	l := labeler.New(nil)
	if got := l.Label(9999); got != "port/9999" {
		t.Fatalf("expected port/9999, got %q", got)
	}
}

func TestCustomRuleSinglePort(t *testing.T) {
	rules := []labeler.Rule{{Low: 8123, High: 8123, Label: "my-app"}}
	l := labeler.New(rules)
	if got := l.Label(8123); got != "my-app" {
		t.Fatalf("expected my-app, got %q", got)
	}
}

func TestCustomRuleRange(t *testing.T) {
	rules := []labeler.Rule{{Low: 9000, High: 9100, Label: "metrics"}}
	l := labeler.New(rules)
	for _, port := range []int{9000, 9050, 9100} {
		if got := l.Label(port); got != "metrics" {
			t.Fatalf("port %d: expected metrics, got %q", port, got)
		}
	}
}

func TestCustomRuleOutsideRange(t *testing.T) {
	rules := []labeler.Rule{{Low: 9000, High: 9100, Label: "metrics"}}
	l := labeler.New(rules)
	if got := l.Label(9101); got != "port/9101" {
		t.Fatalf("expected port/9101, got %q", got)
	}
}

func TestCustomRuleOverridesBuiltin(t *testing.T) {
	rules := []labeler.Rule{{Low: 80, High: 80, Label: "internal-proxy"}}
	l := labeler.New(rules)
	if got := l.Label(80); got != "internal-proxy" {
		t.Fatalf("expected internal-proxy, got %q", got)
	}
}

func TestAddRulePrepends(t *testing.T) {
	l := labeler.New([]labeler.Rule{{Low: 7000, High: 7000, Label: "old"}})
	l.AddRule(labeler.Rule{Low: 7000, High: 7000, Label: "new"})
	if got := l.Label(7000); got != "new" {
		t.Fatalf("expected new, got %q", got)
	}
}
