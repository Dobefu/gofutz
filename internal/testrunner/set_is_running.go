package testrunner

// SetIsRunning sets if tests are running.
func (t *TestRunner) SetIsRunning(running bool) {
	t.mu.Lock()
	t.isRunning = running
	t.mu.Unlock()
}
