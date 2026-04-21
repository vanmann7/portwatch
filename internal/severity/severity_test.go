package severity_test

import (
	"testing"

	"github.com/user/portwatch/internal/severity"
)

func TestCriticalPorts(t *testing.T) {
	classifier := severity.New(nil)
	criticalPorts := []int{21, 22, 23, 3389, 5900, 4444, 1337}
	for _, port := range criticalPorts {
		if got := classifier.Classify(port); got != severity.Critical {
			t.Errorf("port %d: expected Critical, got %s", port, got)
		}
	}
}

func TestWarningPorts(t *testing.T) {
	classifier := severity.New(nil)
	warningPorts := []int{80, 443, 8080, 3306, 5432, 6379, 27017}
	for _, port := range warningPorts {
		if got := classifier.Classify(port); got != severity.Warning {
			t.Errorf("port %d: expected Warning, got %s", port, got)
		}
	}
}

func TestRegisteredRangeIsWarning(t *testing.T) {
	classifier := severity.New(nil)
	// Mid-range registered port not explicitly listed
	if got := classifier.Classify(9090); got != severity.Warning {
		t.Errorf("port 9090: expected Warning, got %s", got)
	}
}

func TestUnknownPortIsInfo(t *testing.T) {
	classifier := severity.New(nil)
	// Port in dynamic/private range above 49151
	if got := classifier.Classify(60000); got != severity.Info {
		t.Errorf("port 60000: expected Info, got %s", got)
	}
}

func TestOverrideTakesPrecedence(t *testing.T) {
	overrides := map[int]severity.Level{
		80:  severity.Info,     // demote HTTP to Info
		9999: severity.Critical, // promote custom port to Critical
	}
	classifier := severity.New(overrides)

	if got := classifier.Classify(80); got != severity.Info {
		t.Errorf("port 80 with override: expected Info, got %s", got)
	}
	if got := classifier.Classify(9999); got != severity.Critical {
		t.Errorf("port 9999 with override: expected Critical, got %s", got)
	}
}

func TestLevelString(t *testing.T) {
	cases := []struct {
		lvl  severity.Level
		want string
	}{
		{severity.Info, "INFO"},
		{severity.Warning, "WARNING"},
		{severity.Critical, "CRITICAL"},
	}
	for _, tc := range cases {
		if got := tc.lvl.String(); got != tc.want {
			t.Errorf("Level.String(): got %q, want %q", got, tc.want)
		}
	}
}

func TestNilOverridesDoNotPanic(t *testing.T) {
	classifier := severity.New(nil)
	_ = classifier.Classify(22) // should not panic
}
