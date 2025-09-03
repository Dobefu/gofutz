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
func (t *TestRunner) RunAllTests(testCompleteCallback func(file File) error) {
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

		files := t.ParseCoverageLines(coverageLines)

		for _, file := range files {
			err = testCompleteCallback(file)

			if err != nil {
				slog.Error(err.Error())
			}
		}
	}()
}
