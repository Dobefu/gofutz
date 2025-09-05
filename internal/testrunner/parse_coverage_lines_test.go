package testrunner

import (
	"testing"
)

func TestParseCoverageLines(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		lines []CoverageLine
		files map[string]File
	}{
		{
			name: "success",
			lines: []CoverageLine{
				{
					File:               "test.go",
					StartColumn:        1,
					StartLine:          1,
					EndLine:            1,
					EndColumn:          1,
					NumberOfStatements: 1,
					ExecutionCount:     1,
				},
			},
			files: map[string]File{
				"test.go": {
					Name: "test.go",
					Functions: []Function{
						{
							Name: "test",
							Result: TestResult{
								Coverage: 1,
							},
						},
					},
					Code:            "",
					HighlightedCode: "",
					Status:          TestStatusPending,
					Coverage:        -1,
					CoveredLines:    []Line{},
				},
				"other.go": {
					Name: "other.go",
					Functions: []Function{
						{
							Name: "other",
							Result: TestResult{
								Coverage: -1,
							},
						},
					},
					Code:            "",
					HighlightedCode: "",
					Status:          TestStatusPending,
					Coverage:        -1,
					CoveredLines:    []Line{},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			runner := &TestRunner{files: test.files} // nolint:exhaustruct

			coveragePercentages := make(map[string]map[string]float64)
			files := runner.ParseCoverageLines(test.lines, coveragePercentages)

			if len(files) != len(test.files) {
				t.Fatalf("expected %d files, got: %d", len(test.files), len(files))
			}

			for _, file := range files {
				if file.Name != test.files[file.Name].Name {
					t.Fatalf(
						"expected file name to be %s, got: %s",
						test.files[file.Name].Name,
						file.Name,
					)
				}
			}
		})
	}
}
