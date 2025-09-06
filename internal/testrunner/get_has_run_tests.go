package testrunner

// GetHasRunTests checks if tests have been run.
func (t *TestRunner) GetHasRunTests() bool {
	t.mu.Lock()
	result := t.hasRunTests
	t.mu.Unlock()

	return result
}
