// Package testrunner provides test runner functionality.
package testrunner

import (
	"fmt"
	"log/slog"

	"github.com/Dobefu/gofutz/internal/filewatcher"
)

// TestRunner defines a test runner.
type TestRunner struct {
	files       map[string]File
	hasRunTests bool
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

func (t *TestRunner) handleFileEvent(path, operation string) {
	moduleName := GetModuleName()

	if moduleName != "" {
		path = fmt.Sprintf("%s/%s", moduleName, path)
	}

	switch operation {
	case "CREATE":
		delete(t.files, path)

		t.files[path] = File{
			Name:            path,
			Functions:       []Function{},
			Code:            "",
			HighlightedCode: "",
			Status:          TestStatusPending,
			Coverage:        -1,
			CoveredLines:    []Line{},
		}

		fallthrough

	case "WRITE", "MODIFY":
		functions, code, err := GetFunctionsFromFile(path)

		if err != nil {
			slog.Error(err.Error())

			return
		}

		var status TestStatus = TestStatusPending

		if len(functions) == 0 {
			status = TestStatusNoCodeToCover
		}

		t.files[path] = File{
			Name:            path,
			Functions:       functions,
			Code:            code,
			HighlightedCode: HighlightCode("go", string(code)),
			Status:          status,
			Coverage:        -1,
			CoveredLines:    []Line{},
		}

	case "REMOVE":
		delete(t.files, path)
	}
}
