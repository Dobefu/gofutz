package testrunner

// GetOutput gets the output buffer.
func (t *TestRunner) GetOutput() []string {
	t.mu.Lock()
	result := t.output
	t.mu.Unlock()

	return result
}
