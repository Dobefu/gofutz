package testrunner

import (
	"testing"
)

func TestRunAllTests(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		files    map[string]File
		expected []string
	}{
		{
			name:     "no test files",
			files:    map[string]File{},
			expected: []string{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			runner := &TestRunner{ // nolint:exhaustruct
				files: test.files,
			}

			fullOutput := []string{}

			runner.RunAllTests(func(_ File) error {
				return nil
			}, func(output string) error {
				fullOutput = append(fullOutput, output)

				return nil
			}, func() {})

			if len(fullOutput) != len(test.expected) {
				t.Fatalf(
					"expected %d output lines, got: %d",
					len(test.expected),
					len(fullOutput),
				)
			}

			for i := range fullOutput {
				if fullOutput[i] != test.expected[i] {
					t.Fatalf(
						"expected output to be %s, got: %s",
						test.expected[i],
						fullOutput[i],
					)
				}
			}
		})
	}
}
