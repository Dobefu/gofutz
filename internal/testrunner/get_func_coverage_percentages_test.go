package testrunner

import (
	"os"
	"testing"
)

func TestGetFuncCoveragePercentages(t *testing.T) {
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

			percentages, _, err := runner.GetFuncCoveragePercentages(
				os.DevNull,
			)

			if err != nil {
				t.Fatalf("expected no error, got: %s", err.Error())
			}

			if len(percentages) != 1 {
				t.Fatalf("expected 1 percentage, got: %d", len(percentages))
			}
		})
	}
}

func TestGetFuncCoveragePercentagesErr(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		expected string
	}{
		{
			name:     "error",
			expected: "exit status 2",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
		})

		runner := &TestRunner{} // nolint:exhaustruct

		_, _, err := runner.GetFuncCoveragePercentages(
			"/bogus",
		)

		if err == nil {
			t.Fatalf("expected error, got nil")
		}

		if err.Error() != test.expected {
			t.Fatalf(
				"expected error to be \"%s\", got: \"%s\"",
				test.expected,
				err.Error(),
			)
		}
	}
}
