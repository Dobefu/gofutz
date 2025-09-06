package testrunner

import (
	"testing"
)

func TestSetIsRunning(t *testing.T) {
	t.Parallel()

	runner := &TestRunner{ // nolint:exhaustruct
		isRunning: false,
	}

	runner.SetIsRunning(true)

	if !runner.isRunning {
		t.Fatalf("expected isRunning to be true, got: %t", runner.isRunning)
	}
}
