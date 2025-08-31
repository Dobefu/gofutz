// Package testrunner provides test runner functionality.
package testrunner

import (
	"log/slog"
	"slices"

	"github.com/Dobefu/gofutz/internal/filewatcher"
)

// TestRunner defines a test runner.
type TestRunner struct {
	files []string
	tests map[string]File
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

	filewatcher.AddListener(func(path, operation string) {
		go runner.handleFileEvent(path, operation)
	})

	return runner, nil
}

// GetTests gets the tests.
func (t *TestRunner) GetTests() map[string]File {
	return t.tests
}

func (t *TestRunner) handleFileEvent(path, operation string) {
	switch operation {
	case "CREATE":
		t.files = slices.DeleteFunc(t.files, func(file string) bool {
			return file == path
		})

		delete(t.tests, path)

		t.files = append(t.files, path)
		slices.Sort(t.files)

		fallthrough

	case "WRITE", "MODIFY":
		tests, err := GetTestsFromFile(path)

		if err != nil {
			slog.Error(err.Error())

			return
		}

		t.tests[path] = File{
			Name:  path,
			Tests: tests,
		}

	case "REMOVE":
		t.files = slices.DeleteFunc(t.files, func(file string) bool {
			return file == path
		})

		delete(t.tests, path)
	}
}
