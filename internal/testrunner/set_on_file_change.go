package testrunner

// SetOnFileChange sets the file change callback.
func (t *TestRunner) SetOnFileChange(callback func()) {
	t.mu.Lock()
	t.onFileChange = callback
	t.mu.Unlock()
}
