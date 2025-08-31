// Package testrunner provides test runner functionality.
package testrunner

import (
	"fmt"
	"log/slog"

	"github.com/Dobefu/gofutz/internal/filewatcher"
)

// TestRunner defines a test runner.
type TestRunner struct {
	files []string
	tests []string
}

// NewTestRunner creates a new test runner.
func NewTestRunner(files []string) (*TestRunner, error) {
	tests, err := GetTestsFromFiles(files)

	if err != nil {
		return nil, err
	}

	runner := &TestRunner{
		files: files,
		tests: tests,
	}

	filewatcher.AddListener(runner.handleFileEvent)

	return runner, nil
}

// GetTests gets the tests.
func (t *TestRunner) GetTests() []string {
	return t.tests
}

func (t *TestRunner) handleFileEvent(path, operation string) {
	slog.Info(
		fmt.Sprintf("file event received: %s %s", operation, path),
	)
}
