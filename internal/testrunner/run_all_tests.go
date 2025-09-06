package testrunner

import (
	"context"
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
	outputCallback func(output string) error,
	completionCallback func(),
) {
	t.mu.Lock()

	if t.cancelFunc != nil {
		t.cancelFunc()
	}

	ctx, cancel := context.WithCancel(context.Background())
	t.cancelFunc = cancel
	t.mu.Unlock()

	go t.doRunAllTests(
		ctx,
		testCompleteCallback,
		outputCallback,
		completionCallback,
	)
}

func (t *TestRunner) doRunAllTests(
	ctx context.Context,
	testCompleteCallback func(file File) error,
	outputCallback func(output string) error,
	completionCallback func(),
) {
	coverageFile, err := os.CreateTemp("", "coverage.out")

	if err != nil {
		slog.Error(fmt.Sprintf("could not create coverage file: %s", err.Error()))
		completionCallback()

		return
	}

	defer func() { _ = coverageFile.Close() }()
	defer func() { _ = os.Remove(coverageFile.Name()) }()

	startTime := time.Now()
	_ = outputCallback("Running all tests")

	goPath, err := exec.LookPath("go")

	if err != nil {
		slog.Error("go command not found")

		return
	}

	cmd := exec.CommandContext(ctx, // #nosec G204 - The temp directory is safe.
		filepath.Clean(goPath),
		"test",
		"-json",
		"-coverprofile",
		filepath.Clean(coverageFile.Name()),
		"./...",
	)

	output, err := cmd.CombinedOutput()

	if ctx.Err() == context.Canceled {
		_ = outputCallback("Tests stopped by user")

		t.mu.Lock()
		t.cancelFunc = nil
		t.mu.Unlock()

		completionCallback()

		return
	}

	if err != nil {
		_ = outputCallback(
			fmt.Sprintf(
				"tests failed: %s\n%s",
				err.Error(),
				t.ParseErrorFromOutput(string(output)),
			),
		)

		coverageLines := []CoverageLine{}
		coveragePercentages := make(map[string]map[string]float64)

		failingTests := t.parseFailedTestFilesFromOutput(string(output))

		err = t.sendCallbacks(
			testCompleteCallback,
			completionCallback,
			coverageLines,
			coveragePercentages,
			failingTests,
		)

		if err != nil {
			slog.Error(
				fmt.Sprintf(
					"could not send callbacks after test failure: %s",
					err.Error(),
				),
			)
		}

		t.mu.Lock()
		t.cancelFunc = nil
		t.mu.Unlock()

		return
	}

	_ = outputCallback(string(output))
	_ = outputCallback(
		fmt.Sprintf("tests completed in %s", time.Since(startTime)),
	)

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
		nil,
	)

	if err != nil {
		slog.Error(fmt.Sprintf("could not send callbacks: %s", err.Error()))
	}

	t.mu.Lock()
	t.cancelFunc = nil
	t.mu.Unlock()
}

func (t *TestRunner) sendCallbacks(
	testCompleteCallback func(file File) error,
	completionCallback func(),
	coverageLines []CoverageLine,
	coveragePercentages map[string]map[string]float64,
	failingTests map[string]bool,
) error {
	files := t.ParseCoverageLines(coverageLines, coveragePercentages)

	if len(failingTests) > 0 {
		t.coverage = -1
	}

	for i := range files {
		if failingTests != nil && failingTests[files[i].Name] {
			files[i].Status = TestStatusFailed
		}

		if len(failingTests) > 0 {
			files[i].Coverage = -1
			files[i].CoveredLines = []Line{}

			for j := range files[i].Functions {
				files[i].Functions[j].Result.Coverage = -1
			}
		}

		t.mu.Lock()
		t.files[files[i].Name] = files[i]
		t.mu.Unlock()
	}

	completionCallback()

	for _, file := range files {
		err := testCompleteCallback(file)

		if err != nil {
			slog.Error(err.Error())
		}
	}

	return nil
}
