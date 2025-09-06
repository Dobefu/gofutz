package testrunner

// GetCoverage gets the overall coverage percentage from the last test run.
func (t *TestRunner) GetCoverage() float64 {
	t.mu.Lock()
	result := t.coverage
	t.mu.Unlock()

	return result
}
