package testrunner

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// CoverageLine defines a coverage line.
type CoverageLine struct {
	File               string `json:"file"`
	StartLine          int    `json:"startLine"`
	StartColumn        int    `json:"startColumn"`
	EndLine            int    `json:"endLine"`
	EndColumn          int    `json:"endColumn"`
	ExecutionCount     int    `json:"executionCount"`
	NumberOfStatements int    `json:"numberOfStatements"`
}

// ParseCoverage parses the coverage report.
func (t *TestRunner) ParseCoverage(coverageFile string) ([]CoverageLine, error) {
	coverage, err := os.ReadFile(filepath.Clean(coverageFile))

	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(coverage), "\n")
	coverageLines := []CoverageLine{}

	for _, line := range lines {
		if line == "" || strings.HasPrefix(line, "mode:") {
			continue
		}

		startLine, startColumn, err := getCoverageStartLineAndColumn(line)

		if err != nil {
			return nil, fmt.Errorf("coverage file not valid: %s", err.Error())
		}

		endLine, endColumn, err := getCoverageEndLineAndColumn(line)

		if err != nil {
			return nil, fmt.Errorf("coverage file not valid: %s", err.Error())
		}

		executionCount, err := getCoverageExecutionCount(line)

		if err != nil {
			return nil, fmt.Errorf("coverage file not valid: %s", err.Error())
		}

		numberOfStatements, err := getCoverageNumberOfStatements(line)

		if err != nil {
			return nil, fmt.Errorf("coverage file not valid: %s", err.Error())
		}

		coverageLines = append(coverageLines, CoverageLine{
			File:               t.getCoverageFileName(line),
			StartLine:          startLine,
			EndLine:            endLine,
			EndColumn:          endColumn,
			StartColumn:        startColumn,
			ExecutionCount:     executionCount,
			NumberOfStatements: numberOfStatements,
		})
	}

	return coverageLines, nil
}

func (t *TestRunner) getCoverageFileName(line string) string {
	return line[:strings.Index(line, ":")]
}

func getCoverageStartLineAndColumn(line string) (int, int, error) {
	startLineAndColumn := regexp.MustCompile(
		`(\d+)\.(\d+)`,
	).FindAllString(line, -1)[0]

	startLine, err := strconv.Atoi(
		startLineAndColumn[:strings.Index(startLineAndColumn, ".")],
	)

	if err != nil {
		return 0, 0, err
	}

	startColumn, err := strconv.Atoi(
		startLineAndColumn[strings.Index(startLineAndColumn, ".")+1:],
	)

	if err != nil {
		return 0, 0, err
	}

	return startLine, startColumn, nil
}

func getCoverageEndLineAndColumn(line string) (int, int, error) {
	endLineAndColumn := regexp.MustCompile(
		`(\d+)\.(\d+)`,
	).FindAllString(line, -1)[1]

	endLine, err := strconv.Atoi(
		endLineAndColumn[:strings.Index(endLineAndColumn, ".")],
	)

	if err != nil {
		return 0, 0, err
	}

	endColumn, err := strconv.Atoi(
		endLineAndColumn[strings.Index(endLineAndColumn, ".")+1:],
	)

	if err != nil {
		return 0, 0, err
	}

	return endLine, endColumn, nil
}

func getCoverageExecutionCount(line string) (int, error) {
	return strconv.Atoi(regexp.MustCompile(
		`\s(\d+)`,
	).FindAllStringSubmatch(line, -1)[0][1])
}

func getCoverageNumberOfStatements(line string) (int, error) {
	return strconv.Atoi(regexp.MustCompile(
		`\s(\d+)`,
	).FindAllStringSubmatch(line, -1)[1][1])
}
