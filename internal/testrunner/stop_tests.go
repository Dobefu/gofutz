package testrunner

// StopTests stops the currently running tests.
func (t *TestRunner) StopTests() {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.cancelFunc != nil {
		t.cancelFunc()
		t.cancelFunc = nil
	}
}
