package testrunner

import (
	"testing"
)

func TestSetCoverage(t *testing.T) {
	t.Parallel()

	runner := &TestRunner{ // nolint:exhaustruct
		coverage: -1,
	}

	runner.SetCoverage(100)

	if runner.coverage != 100 {
		t.Fatalf("expected coverage to be 100, got: %f", runner.coverage)
	}
}
