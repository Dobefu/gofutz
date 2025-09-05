package testrunner

import (
	"os"
	"testing"
)

func TestParseCoverage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		coverageFile  string
		expectedLines []CoverageLine
	}{
		{
			name:          "nonexistent file",
			coverageFile:  os.DevNull,
			expectedLines: []CoverageLine{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			runner := &TestRunner{} // nolint:exhaustruct

			lines, err := runner.ParseCoverage(test.coverageFile)

			if err != nil {
				t.Fatalf("expected no error, got: %s", err.Error())
			}

			if len(lines) != len(test.expectedLines) {
				t.Fatalf("expected %d lines, got: %d", len(test.expectedLines), len(lines))
			}

			for i, line := range lines {
				if line.File != test.expectedLines[i].File {
					t.Fatalf("expected file to be %s, got: %s", test.expectedLines[i].File, line.File)
				}
			}
		})
	}
}

func TestGetCoverageStartLineAndColumn(t *testing.T) {
	t.Parallel()

	startLine, startColumn, err := getCoverageStartLineAndColumn("test.go:1.2:0.0 0 0")

	if err != nil {
		t.Errorf("expected no error, got: %s", err.Error())
	}

	if startLine != 1 {
		t.Errorf("expected start line to be 1, got: %d", startLine)
	}

	if startColumn != 2 {
		t.Errorf("expected start column to be 2, got: %d", startColumn)
	}
}

func TestGetCoverageEndLineAndColumn(t *testing.T) {
	t.Parallel()

	endLine, endColumn, err := getCoverageEndLineAndColumn("test.go:0.0:1.2 0 0")

	if err != nil {
		t.Errorf("expected no error, got: %s", err.Error())
	}

	if endLine != 1 {
		t.Errorf("expected end line to be 1, got: %d", endLine)
	}

	if endColumn != 2 {
		t.Errorf("expected end column to be 2, got: %d", endColumn)
	}
}

func TestGetCoverageNumberOfStatements(t *testing.T) {
	t.Parallel()

	numberOfStatements, err := getCoverageNumberOfStatements("test.go:0.0:0.0 1 0")

	if err != nil {
		t.Errorf("expected no error, got: %s", err.Error())
	}

	if numberOfStatements != 1 {
		t.Errorf("expected number of statements to be 1, got: %d", numberOfStatements)
	}
}

func TestGetCoverageExecutionCount(t *testing.T) {
	t.Parallel()

	numberOfStatements, err := getCoverageExecutionCount("test.go:0.0:0.0 0 1")

	if err != nil {
		t.Errorf("expected no error, got: %s", err.Error())
	}

	if numberOfStatements != 1 {
		t.Errorf("expected number of statements to be 1, got: %d", numberOfStatements)
	}
}
