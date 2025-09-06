package testrunner

import (
	"testing"
)

func TestParseErrorFromOutput(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		output   string
		expected string
	}{
		{
			name:     "unexpected output",
			output:   `{"time":"2021-01-01T00:00:00Z","action":"fail","package":"/path/to/test","test":"TestSomething","output":"expected: 1, got: 2","elapsed":0}`,
			expected: "expected: 1, got: 2",
		},
		{
			name:     "empty output",
			output:   "",
			expected: "",
		},
		{
			name:     "invalid JSON",
			output:   "{",
			expected: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			runner := &TestRunner{} // nolint:exhaustruct
			output := runner.ParseErrorFromOutput(test.output)

			if output != test.expected {
				t.Fatalf("expected \"%s\", got \"%s\"", test.expected, output)
			}
		})
	}
}
