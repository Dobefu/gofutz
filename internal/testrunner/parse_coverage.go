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
	NumberOfStatements int    `json:"numberOfStatements"`
	ExecutionCount     int    `json:"executionCount"`
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
			return nil, err
		}

		endLine, endColumn, err := getCoverageEndLineAndColumn(line)

		if err != nil {
			return nil, err
		}

		numberOfStatements, err := getCoverageNumberOfStatements(line)

		if err != nil {
			return nil, err
		}

		executionCount, err := getCoverageExecutionCount(line)

		if err != nil {
			return nil, err
		}

		coverageLines = append(coverageLines, CoverageLine{
			File:               t.getCoverageFileName(line),
			StartLine:          startLine,
			EndLine:            endLine,
			EndColumn:          endColumn,
			StartColumn:        startColumn,
			NumberOfStatements: numberOfStatements,
			ExecutionCount:     executionCount,
		})
	}

	return coverageLines, nil
}

func (t *TestRunner) getCoverageFileName(line string) string {
	return line[:strings.Index(line, ":")]
}

func getCoverageStartLineAndColumn(line string) (int, int, error) {
	startLineAndColumnMatches := regexp.MustCompile(
		`(\d+)\.(\d+)`,
	).FindAllString(line, -1)

	if len(startLineAndColumnMatches) < 1 {
		return 0, 0, fmt.Errorf("missing start line and column: %s", line)
	}

	startLineAndColumn := startLineAndColumnMatches[0]

	// This cannot fail, since we explicitly check for digits.
	startLine, _ := strconv.Atoi(
		startLineAndColumn[:strings.Index(startLineAndColumn, ".")],
	)

	// This cannot fail, since we explicitly check for digits.
	startColumn, _ := strconv.Atoi(
		startLineAndColumn[strings.Index(startLineAndColumn, ".")+1:],
	)

	return startLine, startColumn, nil
}

func getCoverageEndLineAndColumn(line string) (int, int, error) {
	endLineAndColumnMatches := regexp.MustCompile(
		`(\d+)\.(\d+)`,
	).FindAllString(line, -1)

	if len(endLineAndColumnMatches) < 2 {
		return 0, 0, fmt.Errorf("missing end line and column: %s", line)
	}

	endLineAndColumn := endLineAndColumnMatches[1]

	// This cannot fail, since we explicitly check for digits.
	endLine, _ := strconv.Atoi(
		endLineAndColumn[:strings.Index(endLineAndColumn, ".")],
	)

	// This cannot fail, since we explicitly check for digits.
	endColumn, _ := strconv.Atoi(
		endLineAndColumn[strings.Index(endLineAndColumn, ".")+1:],
	)

	return endLine, endColumn, nil
}

func getCoverageNumberOfStatements(line string) (int, error) {
	numberOfStatementsMatches := regexp.MustCompile(
		`\s(\d+)`,
	).FindAllStringSubmatch(line, -1)

	if len(numberOfStatementsMatches) < 1 {
		return 0, fmt.Errorf("missing number of statements: %s", line)
	}

	return strconv.Atoi(numberOfStatementsMatches[0][1])
}

func getCoverageExecutionCount(line string) (int, error) {
	executionCountMatches := regexp.MustCompile(
		`\s(\d+)`,
	).FindAllStringSubmatch(line, -1)

	if len(executionCountMatches) < 2 {
		return 0, fmt.Errorf("missing execution count: %s", line)
	}

	return strconv.Atoi(executionCountMatches[1][1])
}
