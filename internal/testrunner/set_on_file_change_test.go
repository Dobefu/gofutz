package testrunner

import (
	"testing"
)

func TestSetOnFileChange(t *testing.T) {
	t.Parallel()

	runner := &TestRunner{ // nolint:exhaustruct
		onFileChange: nil,
	}

	runner.SetOnFileChange(func() {})

	if runner.onFileChange == nil {
		t.Fatalf("expected onFileChange to be set, got: nil")
	}
}
