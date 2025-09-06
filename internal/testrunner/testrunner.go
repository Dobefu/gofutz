// Package testrunner provides test runner functionality.
package testrunner

import (
	"context"
	"sync"
	"time"

	"github.com/Dobefu/gofutz/internal/filewatcher"
)

// TestRunner defines a test runner.
type TestRunner struct {
	files            map[string]File
	hasRunTests      bool
	coverage         float64
	output           []string
	isRunning        bool
	debounceFiles    map[string]*time.Timer
	debounceDuration time.Duration
	mu               sync.Mutex
	onFileChange     func()
	cancelFunc       context.CancelFunc
}

// NewTestRunner creates a new test runner.
func NewTestRunner(
	files []string,
	fw *filewatcher.FileWatcher,
	onFileChange func(),
) (*TestRunner, error) {
	tests, err := GetFunctionsFromFiles(files)

	if err != nil {
		return nil, err
	}

	runner := &TestRunner{
		files:            tests,
		hasRunTests:      false,
		coverage:         -1,
		output:           []string{},
		isRunning:        false,
		debounceFiles:    make(map[string]*time.Timer),
		debounceDuration: 100 * time.Millisecond,
		mu:               sync.Mutex{},
		onFileChange:     onFileChange,
		cancelFunc:       nil,
	}

	fw.AddListener(func(path, operation string) {
		go runner.handleFileEvent(path, operation)
	})

	return runner, nil
}
