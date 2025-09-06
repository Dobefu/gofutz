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

func TestSendCallbacks(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		files    map[string]File
		expected map[string]bool
	}{
		{
			name: "success",
			files: map[string]File{
				"test.go": { // nolint:exhaustruct
					Name: "test.go",
				},
			},
			expected: map[string]bool{
				"test.go": true,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			runner := &TestRunner{ // nolint:exhaustruct
				files: test.files,
			}

			err := runner.sendCallbacks(
				func(_ File) error {
					return nil
				},
				func() {},
				[]CoverageLine{},
				map[string]map[string]float64{},
				test.expected,
			)

			if err != nil {
				t.Fatalf("expected no error, got: %s", err.Error())
			}

			if len(runner.files) != len(test.expected) {
				t.Fatalf(
					"expected %d files, got: %d",
					len(test.expected),
					len(runner.files),
				)
			}

			for file := range runner.files {
				if runner.files[file].Status != TestStatusFailed {
					t.Fatalf(
						"expected file status to be \"%T\", got: \"%T\"",
						TestStatusFailed,
						runner.files[file].Status,
					)
				}
			}
		})
	}
}
