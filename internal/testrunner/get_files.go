package testrunner

import (
	"maps"
)

// GetFiles gets the files.
func (t *TestRunner) GetFiles() map[string]File {
	t.mu.Lock()
	defer t.mu.Unlock()

	files := make(map[string]File)
	maps.Copy(files, t.files)

	return files
}
