package testrunner

import (
	"testing"
)

func TestSetHasRunTests(t *testing.T) {
	t.Parallel()

	runner := &TestRunner{ // nolint:exhaustruct
		hasRunTests: false,
	}

	runner.SetHasRunTests(true)

	if !runner.hasRunTests {
		t.Fatalf("expected hasRunTests to be true, got: %t", runner.hasRunTests)
	}
}
