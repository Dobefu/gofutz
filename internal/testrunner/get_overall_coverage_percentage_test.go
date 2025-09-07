package testrunner

import (
	"testing"
)

func TestGetOverallCoveragePercentage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		output   []byte
		expected float64
	}{
		{
			name:     "success",
			output:   []byte("total: 100.0%"),
			expected: 100.0,
		},
		{
			name:     "empty output",
			output:   []byte(""),
			expected: -1,
		},
		{
			name:     "missing total",
			output:   []byte("total:"),
			expected: -1,
		},
		{
			name:     "invalid output",
			output:   []byte("total: invalid"),
			expected: -1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			coverage := GetOverallCoveragePercentage(test.output)

			if coverage != test.expected {
				t.Fatalf("expected coverage to be %f, got: %f", test.expected, coverage)
			}
		})
	}
}
