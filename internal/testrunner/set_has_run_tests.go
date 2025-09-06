package testrunner

// SetHasRunTests sets if tests have been run.
func (t *TestRunner) SetHasRunTests(hasRunTests bool) {
	t.mu.Lock()
	t.hasRunTests = hasRunTests
	t.mu.Unlock()
}
