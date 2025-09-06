package testrunner

// GetIsRunning checks if tests are currently running.
func (t *TestRunner) GetIsRunning() bool {
	return t.isRunning
}
