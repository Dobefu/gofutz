package testrunner

// AddOutput adds output lines to the output buffer.
func (t *TestRunner) AddOutput(output []string) []string {
	t.output = append(t.output, output...)

	maxLines := 200

	if len(t.output) > maxLines {
		t.output = t.output[len(t.output)-maxLines:]
	}

	return t.output
}
