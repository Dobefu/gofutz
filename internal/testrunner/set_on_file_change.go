package testrunner

// SetOnFileChange sets the file change callback.
func (t *TestRunner) SetOnFileChange(callback func()) {
	t.onFileChange = callback
}
