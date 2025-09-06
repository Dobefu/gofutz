// Package testrunner provides test runner functionality.
package testrunner

import (
	"sync"
	"time"

	"github.com/Dobefu/gofutz/internal/filewatcher"
)

// TestRunner defines a test runner.
type TestRunner struct {
	files         map[string]File
	hasRunTests   bool
	coverage      float64
	output        []string
	isRunning     bool
	debounceFiles map[string]*time.Timer
	mu            sync.Mutex
	onFileChange  func()
}

// NewTestRunner creates a new test runner.
func NewTestRunner(files []string, onFileChange func()) (*TestRunner, error) {
	tests, err := GetFunctionsFromFiles(files)

	if err != nil {
		return nil, err
	}

	runner := &TestRunner{
		files:         tests,
		hasRunTests:   false,
		coverage:      -1,
		output:        []string{},
		isRunning:     false,
		debounceFiles: make(map[string]*time.Timer),
		mu:            sync.Mutex{},
		onFileChange:  onFileChange,
	}

	filewatcher.AddListener(func(path, operation string) {
		go runner.handleFileEvent(path, operation)
	})

	return runner, nil
}
