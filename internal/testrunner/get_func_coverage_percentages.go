package testrunner

import (
	"log/slog"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// GetFuncCoveragePercentages gets coverage percentages for each function.
func (t *TestRunner) GetFuncCoveragePercentages(
	coverageFile string,
) (map[string]map[string]float64, error) {
	goPath, err := exec.LookPath("go")

	if err != nil {
		slog.Error("go command not found")

		return nil, err
	}

	cmd := exec.Command( // #nosec G204 - The temp directory is safe.
		filepath.Clean(goPath),
		"tool",
		"cover",
		"-func",
		filepath.Clean(coverageFile),
	)

	output, err := cmd.CombinedOutput()

	if err != nil {
		return nil, err
	}

	percentages := make(map[string]map[string]float64)
	lines := strings.SplitSeq(string(output), "\n")

	for line := range lines {
		if !strings.Contains(line, "%") || !strings.Contains(line, ":") {
			continue
		}

		parts := strings.Fields(line)

		if len(parts) < 2 {
			continue
		}

		fileName := t.getCoverageFileName(line)
		testName := parts[1]

		testPercentage := parts[len(parts)-1]
		testPercentage = strings.TrimSuffix(testPercentage, "%")

		percentage, err := strconv.ParseFloat(testPercentage, 64)

		if err != nil {
			continue
		}

		_, hasFile := percentages[fileName]

		if !hasFile {
			percentages[fileName] = make(map[string]float64)
		}

		percentages[fileName][testName] = percentage
	}

	return percentages, nil
}
