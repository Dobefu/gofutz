package testrunner

// SetCoverage sets the overall coverage percentage.
func (t *TestRunner) SetCoverage(coverage float64) {
	t.mu.Lock()
	t.coverage = coverage
	t.mu.Unlock()
}
