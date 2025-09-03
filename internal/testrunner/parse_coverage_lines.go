package testrunner

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

	for fileName, lines := range coverage {
		coveragePercentage := coveragePercentages[fileName]
		file, hasFile := t.files[fileName]

		if !hasFile {
			continue
		}

		for i, function := range file.Functions {
			function.Result.Coverage = coveragePercentage[function.Name]
			file.Functions[i] = function
		}

		var numCoveredStatements int
		var totalStatements int

		for _, line := range coverage[fileName] {
			totalStatements += line.NumberOfStatements

			if line.ExecutionCount > 0 {
				numCoveredStatements += line.NumberOfStatements
			}
		}

		file.Coverage = (float64(numCoveredStatements) / float64(totalStatements)) * 100
		file.CoveredLines = lines

		t.files[fileName] = file

		files = append(files, file)
	}

	return files
}
