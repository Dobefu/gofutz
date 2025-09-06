package testrunner

import (
	"testing"
)

func TestGetHasRunTests(t *testing.T) {
	t.Parallel()

	runner := &TestRunner{ // nolint:exhaustruct
		hasRunTests: true,
	}

	hasRunTests := runner.GetHasRunTests()

	if !hasRunTests {
		t.Fatalf("expected hasRunTests to be true, got: %t", hasRunTests)
	}
}
