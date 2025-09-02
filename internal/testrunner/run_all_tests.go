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
func (t *TestRunner) RunAllTests(testCompleteCallback func(test Test) error) {
	go func() {
		coverageFile := filepath.Join(os.TempDir(), "coverage.out")
		defer func() { _ = os.Remove(coverageFile) }()

		startTime := time.Now()
		slog.Info("Running all tests")

		goPath, err := exec.LookPath("go")

		if err != nil {
			slog.Error("go command not found")

			return
		}

		cmd := exec.Command( // #nosec G204 -- The temp directory is safe.
			filepath.Clean(goPath),
			"test",
			"-v",
			"-coverprofile",
			filepath.Clean(coverageFile),
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

		coverageLines, err := t.ParseCoverage(coverageFile)

		if err != nil {
			slog.Error(fmt.Sprintf("could not parse coverage: %s", err.Error()))
		}

		tests := t.ParseCoverageLines(coverageLines)

		slog.Info(fmt.Sprintf("coverage: %v", tests))

		for _, test := range tests {
			err = testCompleteCallback(Test{
				Name:   test.Name,
				Result: test.Result,
			})

			if err != nil {
				slog.Error(err.Error())
			}
		}
	}()
}
