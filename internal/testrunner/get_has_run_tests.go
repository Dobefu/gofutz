package testrunner

// GetHasRunTests checks if tests have been run.
func (t *TestRunner) GetHasRunTests() bool {
	return t.hasRunTests
}
