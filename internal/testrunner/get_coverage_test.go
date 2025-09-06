package testrunner

import (
	"testing"
)

func TestGetCoverage(t *testing.T) {
	t.Parallel()

	runner := &TestRunner{ // nolint:exhaustruct
		coverage: -1,
	}

	coverage := runner.GetCoverage()

	if coverage != -1 {
		t.Fatalf("expected coverage to be -1, got: %f", coverage)
	}
}
