package testrunner

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func createCoverageFile(
	t *testing.T,
	name string,
	coverageString string,
) (string, func()) {
	t.Helper()

	if coverageString == "" {
		return "/bogus", func() {}
	}

	coverageFile := filepath.Join(
		os.TempDir(),
		fmt.Sprintf("%s.coverage", name),
	)

	err := os.WriteFile(coverageFile, []byte(coverageString), 0600)

	if err != nil {
		t.Fatalf("expected no error, got: %s", err.Error())
	}

	return coverageFile, func() { _ = os.Remove(coverageFile) }
}

func TestParseCoverage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		coverageString string
		expectedLines  []CoverageLine
	}{
		{
			name:           "success",
			coverageString: "mode:set\ntest.go:1.2:3.4 5 6",
			expectedLines: []CoverageLine{
				{
					File:               "test.go",
					StartLine:          1,
					StartColumn:        2,
					EndLine:            3,
					EndColumn:          4,
					NumberOfStatements: 5,
					ExecutionCount:     6,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			coverageFile, cleanup := createCoverageFile(
				t,
				test.name,
				test.coverageString,
			)

			defer cleanup()

			runner := &TestRunner{} // nolint:exhaustruct

			lines, err := runner.ParseCoverage(coverageFile)

			if err != nil {
				t.Fatalf("expected no error, got: %s", err.Error())
			}

			if len(lines) != len(test.expectedLines) {
				t.Fatalf("expected %d lines, got: %d", len(test.expectedLines), len(lines))
			}

			linesJSON, err := json.Marshal(lines)

			if err != nil {
				t.Fatalf("expected no error, got: %s", err.Error())
			}

			expectedLinesJSON, _ := json.Marshal(test.expectedLines)

			if string(linesJSON) != string(expectedLinesJSON) {
				t.Fatalf("expected \"%s\", got \"%s\"", expectedLinesJSON, string(linesJSON))
			}
		})
	}
}

func TestParseCoverageErr(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		coverageString string
		expected       string
	}{
		{
			name:           "no coverage file",
			coverageString: "",
			expected:       "no such file or directory",
		},
		{
			name:           "invalid coverage file - missing start line",
			coverageString: "mode:set\ntest.go:1",
			expected:       "invalid coverage line:",
		},
		{
			name:           "invalid coverage file - missing end column",
			coverageString: "mode:set\ntest.go:1.2:3",
			expected:       "invalid coverage line:",
		},
		{
			name:           "invalid coverage file - missing number of statements",
			coverageString: "mode:set\ntest.go:1.2:3.4 5",
			expected:       "invalid coverage line:",
		},
		{
			name:           "invalid coverage file - missing execution count",
			coverageString: "mode:set\ntest.go:1.2:3.4",
			expected:       "invalid coverage line:",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			coverageFile, cleanup := createCoverageFile(
				t,
				test.name,
				test.coverageString,
			)

			defer cleanup()

			runner := &TestRunner{} // nolint:exhaustruct

			_, err := runner.ParseCoverage(coverageFile)

			if err == nil {
				t.Fatalf("expected error, got nil")
			}

			if !strings.Contains(err.Error(), test.expected) {
				t.Fatalf(
					"expected error to contain \"%s\", got \"%s\"",
					test.expected,
					err.Error(),
				)
			}
		})
	}
}
