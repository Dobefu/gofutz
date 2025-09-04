// Package testrunner provides test runner functionality.
package testrunner

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/Dobefu/gofutz/internal/filewatcher"
)

// TestRunner defines a test runner.
type TestRunner struct {
	files       map[string]File
	hasRunTests bool
	coverage    float64
	isRunning   bool
}

// NewTestRunner creates a new test runner.
func NewTestRunner(files []string) (*TestRunner, error) {
	tests, err := GetFunctionsFromFiles(files)

	if err != nil {
		return nil, err
	}

	runner := &TestRunner{
		files:       tests,
		hasRunTests: false,
		coverage:    -1,
		isRunning:   false,
	}

	filewatcher.AddListener(func(path, operation string) {
		go runner.handleFileEvent(path, operation)
	})

	return runner, nil
}

// GetFiles gets the files.
func (t *TestRunner) GetFiles() map[string]File {
	return t.files
}

// HasRunTests checks if tests have been run.
func (t *TestRunner) HasRunTests() bool {
	return t.hasRunTests
}

// SetHasRunTests sets if tests have been run.
func (t *TestRunner) SetHasRunTests(hasRunTests bool) {
	t.hasRunTests = hasRunTests
}

// GetCoverage returns the overall coverage percentage from the last test run.
func (t *TestRunner) GetCoverage() float64 {
	return t.coverage
}

// SetCoverage sets the overall coverage percentage.
func (t *TestRunner) SetCoverage(coverage float64) {
	t.coverage = coverage
}

// IsRunning returns true if tests are currently running.
func (t *TestRunner) IsRunning() bool {
	return t.isRunning
}

// SetRunning sets the running state.
func (t *TestRunner) SetRunning(running bool) {
	t.isRunning = running
}

func (t *TestRunner) handleFileEvent(path, operation string) {
	// fmt.Println(path)
	cwd, err := os.Getwd()

	if err != nil {
		slog.Error(err.Error())

		return
	}

	moduleName := GetModuleName()
	modulePath := fmt.Sprintf("%s%s", moduleName, strings.TrimPrefix(path, cwd))

	switch operation {
	case "CREATE":
		delete(t.files, modulePath)

		t.files[modulePath] = File{
			Name:            modulePath,
			Functions:       []Function{},
			Code:            "",
			HighlightedCode: "",
			Status:          TestStatusPending,
			Coverage:        -1,
			CoveredLines:    []Line{},
		}

		fallthrough

	case "WRITE", "MODIFY", "RENAME":
		functions, code, err := GetFunctionsFromFile(path)

		if err != nil {
			slog.Error(err.Error())

			return
		}

		var status TestStatus = TestStatusPending

		if len(functions) == 0 {
			status = TestStatusNoCodeToCover
		}

		t.files[modulePath] = File{
			Name:            modulePath,
			Functions:       functions,
			Code:            code,
			HighlightedCode: HighlightCode("go", string(code)),
			Status:          status,
			Coverage:        -1,
			CoveredLines:    []Line{},
		}

	case "REMOVE":
		delete(t.files, modulePath)
	}
}
