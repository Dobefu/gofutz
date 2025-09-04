package testrunner

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// RunAllTests runs all tests.
func (t *TestRunner) RunAllTests(
	testCompleteCallback func(file File) error,
	completionCallback func() error,
) {
	go func() {
		coverageFile, err := os.CreateTemp("", "coverage.out")

		if err != nil {
			slog.Error(fmt.Sprintf("could not create coverage file: %s", err.Error()))

			return
		}

		defer func() { _ = coverageFile.Close() }()
		defer func() { _ = os.Remove(coverageFile.Name()) }()

		startTime := time.Now()
		slog.Info("Running all tests")

		goPath, err := exec.LookPath("go")

		if err != nil {
			slog.Error("go command not found")

			return
		}

		cmd := exec.Command( // #nosec G204 - The temp directory is safe.
			filepath.Clean(goPath),
			"test",
			"-v",
			"-coverprofile",
			filepath.Clean(coverageFile.Name()),
			"./...",
		)

		output, err := cmd.CombinedOutput()

		if err != nil {
			slog.Error(
				fmt.Sprintf("tests failed: %s: %s", err.Error(), string(output)),
			)

			return
		}

		slog.Info(fmt.Sprintf("tests completed in %s", time.Since(startTime)))

		coverageLines, err := t.ParseCoverage(coverageFile.Name())

		if err != nil {
			slog.Error(fmt.Sprintf("could not parse coverage: %s", err.Error()))
		}

		coveragePercentages, overallCoverage, err := t.GetFuncCoveragePercentages(
			coverageFile.Name(),
		)

		if err != nil {
			slog.Error(
				fmt.Sprintf("could not get coverage percentages: %s", err.Error()),
			)

			return
		}

		t.SetCoverage(overallCoverage)

		err = t.sendCallbacks(
			testCompleteCallback,
			completionCallback,
			coverageLines,
			coveragePercentages,
		)

		if err != nil {
			slog.Error(fmt.Sprintf("could not send callbacks: %s", err.Error()))
		}
	}()
}

func (t *TestRunner) sendCallbacks(
	testCompleteCallback func(file File) error,
	completionCallback func() error,
	coverageLines []CoverageLine,
	coveragePercentages map[string]map[string]float64,
) error {
	files := t.ParseCoverageLines(coverageLines, coveragePercentages)

	for _, file := range files {
		err := testCompleteCallback(file)

		if err != nil {
			slog.Error(err.Error())
		}
	}

	err := completionCallback()

	if err != nil {
		slog.Error(fmt.Sprintf("could not send completion update: %s", err.Error()))
	}

	return nil
}
