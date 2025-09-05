package testrunner

import (
	"testing"
)

func TestParseFailedTestFilesFromOutput(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		output   string
		expected map[string]bool
	}{
		{
			name:   "success",
			output: `{"time":"2021-01-01T00:00:00Z","action":"fail","package":"path/to/test","test":"Test","output":"","elapsed":0}`,
			expected: map[string]bool{
				"/path/to/test.go": true,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			runner := &TestRunner{ // nolint:exhaustruct
				files: map[string]File{
					"test.go": { // nolint:exhaustruct
						Name: "test.go",
					},
					"/path/to/test.go": { // nolint:exhaustruct
						Name: "/path/to/test.go",
					},
				},
			}

			files := runner.parseFailedTestFilesFromOutput(test.output)

			if len(files) != len(test.expected) {
				t.Fatalf("expected %d file(s), got: %d", len(test.expected), len(files))
			}

			for file := range files {
				if files[file] != test.expected[file] {
					t.Fatalf(
						"expected file %s to be %t, got: %t",
						file,
						test.expected[file],
						files[file],
					)
				}
			}
		})
	}
}
