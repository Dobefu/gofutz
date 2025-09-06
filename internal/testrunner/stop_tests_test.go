package testrunner

import (
	"testing"
)

func TestStopTests(t *testing.T) {
	t.Parallel()

	runner := &TestRunner{ // nolint:exhaustruct
		cancelFunc: func() {},
	}

	runner.StopTests()

	if runner.cancelFunc != nil {
		t.Fatalf("expected cancelFunc to be nil, got: %v", runner.cancelFunc)
	}
}
