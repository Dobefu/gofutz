package testrunner

import (
	"maps"
)

// ParseCoverageLines parses the coverage lines.
func (t *TestRunner) ParseCoverageLines(
	coverageLines []CoverageLine,
	coveragePercentages map[string]map[string]float64,
) []File {
	coverage := make(map[string][]Line)

	for _, line := range coverageLines {
		fileName := line.File

		coverage[fileName] = append(coverage[fileName], Line{
			Number:             line.StartLine,
			StartLine:          line.StartLine,
			StartColumn:        line.StartColumn,
			EndLine:            line.EndLine,
			EndColumn:          line.EndColumn,
			ExecutionCount:     line.ExecutionCount,
			NumberOfStatements: line.NumberOfStatements,
		})
	}

	files := []File{}

	t.mu.Lock()
	newFiles := make(map[string]File)
	maps.Copy(newFiles, t.files)
	t.mu.Unlock()

	for fileName, file := range newFiles {
		lines, hasCoverage := coverage[fileName]
		coveragePercentage := coveragePercentages[fileName]

		if !hasCoverage {
			file.Coverage = -1

			if file.Status != TestStatusNoCodeToCover {
				file.Status = TestStatusPending
				file.Coverage = 0
			}

			file.CoveredLines = []Line{}

			for i, function := range file.Functions {
				function.Result.Coverage = 0
				file.Functions[i] = function
			}

			t.mu.Lock()
			t.files[fileName] = file
			t.mu.Unlock()
			files = append(files, file)

			continue
		}

		for i, function := range file.Functions {
			function.Result.Coverage = coveragePercentage[function.Name]
			file.Functions[i] = function
		}

		file.Coverage = getFileCoverage(lines)

		file.CoveredLines = lines

		if file.Status != TestStatusNoCodeToCover {
			file.Status = TestStatusPassed
		}

		t.mu.Lock()
		t.files[fileName] = file
		t.mu.Unlock()

		files = append(files, file)
	}

	return files
}

func getFileCoverage(lines []Line) float64 {
	var numCoveredStatements int
	var totalStatements int

	for _, line := range lines {
		totalStatements += line.NumberOfStatements

		if line.ExecutionCount > 0 {
			numCoveredStatements += line.NumberOfStatements
		}
	}

	if totalStatements > 0 {
		return (float64(numCoveredStatements) / float64(totalStatements)) * 100
	}

	return 0
}
