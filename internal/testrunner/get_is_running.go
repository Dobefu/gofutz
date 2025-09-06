package testrunner

// GetIsRunning checks if tests are currently running.
func (t *TestRunner) GetIsRunning() bool {
	t.mu.Lock()
	result := t.isRunning
	t.mu.Unlock()

	return result
}
