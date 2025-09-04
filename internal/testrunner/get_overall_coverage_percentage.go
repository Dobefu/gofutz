package testrunner

import (
	"strconv"
	"strings"
)

// GetOverallCoveragePercentage gets the overall coverage percentage.
func GetOverallCoveragePercentage(output []byte) float64 {
	lines := strings.SplitSeq(string(output), "\n")

	for line := range lines {
		if !strings.Contains(line, "total:") {
			continue
		}

		parts := strings.Fields(line)

		if len(parts) < 2 {
			continue
		}

		totalPercentage := strings.TrimSuffix(parts[len(parts)-1], "%")
		coverage, err := strconv.ParseFloat(totalPercentage, 64)

		if err != nil {
			continue
		}

		return coverage
	}

	return -1
}
