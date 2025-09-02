package testrunner

import (
	"strings"
)

// ParseCoverageLines parses the coverage lines.
func (t *TestRunner) ParseCoverageLines(coverageLines []CoverageLine) []Test {
	coverage := make(map[string][]Line)

	for _, line := range coverageLines {
		fileName := line.File

		coverage[fileName] = append(coverage[fileName], Line{
			Number:         line.StartLine,
			ExecutionCount: line.ExecutionCount,
		})
	}

	tests := []Test{}

	for fileName, lines := range coverage {
		testFile, hasTestFile := t.files[fileName]

		if !hasTestFile {
			continue
		}

		var coveragePercentage float64
		numLines := len(strings.Split(testFile.Code, "\n"))
		numCoveredLines := 0

		for _, line := range lines {
			numLines++

			if line.ExecutionCount > 0 {
				numCoveredLines++
			}
		}

		if numLines > 0 {
			coveragePercentage = float64(numCoveredLines) / float64(numLines) * 100
		}

		tests = append(tests, Test{
			Name: fileName,
			Result: TestResult{
				Status:       TestStatusRunning,
				Output:       []string{},
				Coverage:     coveragePercentage,
				CoveredLines: lines,
			},
		})
	}

	return tests
}
