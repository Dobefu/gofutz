package testrunner

import (
	"time"
)

// Close closes the test runner.
func (t *TestRunner) Close() {
	t.mu.Lock()
	defer t.mu.Unlock()

	for _, timer := range t.debounceFiles {
		timer.Stop()
	}

	t.debounceFiles = make(map[string]*time.Timer)
	t.isRunning = false
}
