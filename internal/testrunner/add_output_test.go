package testrunner

import (
	"testing"
)

func TestAddOutput(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
	}{
		{
			name: "success",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			runner := &TestRunner{} // nolint:exhaustruct

			output := runner.AddOutput([]string{"test"})

			if len(output) != 1 {
				t.Fatalf("expected 1 output line, got: %d", len(output))
			}

			if output[0] != "test" {
				t.Fatalf("expected output to be test, got: %s", output[0])
			}

			for range 200 {
				output = runner.AddOutput([]string{"test"})
			}

			if len(output) != 200 {
				t.Fatalf("expected 200 output lines, got: %d", len(output))
			}
		})
	}
}
