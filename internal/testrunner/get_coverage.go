package testrunner

// GetCoverage gets the overall coverage percentage from the last test run.
func (t *TestRunner) GetCoverage() float64 {
	return t.coverage
}
