package testrunner

import (
	"testing"
)

func TestGetIsRunning(t *testing.T) {
	t.Parallel()

	runner := &TestRunner{ // nolint:exhaustruct
		isRunning: true,
	}

	isRunning := runner.GetIsRunning()

	if !isRunning {
		t.Fatalf("expected isRunning to be true, got: %t", isRunning)
	}
}
