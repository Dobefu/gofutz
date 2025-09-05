package testrunner

import (
	"testing"
)

func TestGetOutput(t *testing.T) {
	t.Parallel()

	runner := &TestRunner{} // nolint:exhaustruct

	output := runner.GetOutput()

	if len(output) != 0 {
		t.Fatalf("expected empty output, got: %v", output)
	}

	output = runner.AddOutput([]string{"test"})

	if len(output) != 1 {
		t.Fatalf("expected 1 output line, got: %d", len(output))
	}

	if output[0] != "test" {
		t.Fatalf("expected output to be test, got: %s", output[0])
	}
}
