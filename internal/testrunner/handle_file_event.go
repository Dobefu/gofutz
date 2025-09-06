package testrunner

import (
	"strings"
	"time"
)

func (t *TestRunner) handleFileEvent(path, operation string) {
	t.mu.Lock()
	timer, hasTimer := t.debounceFiles[path]

	if hasTimer {
		timer.Stop()
	}

	t.debounceFiles[path] = time.AfterFunc(t.debounceDuration, func() {
		if !strings.HasSuffix(path, ".go") {
			t.mu.Lock()
			delete(t.debounceFiles, path)
			t.mu.Unlock()

			return
		}

		t.processFileEvent(path, operation)

		t.mu.Lock()
		shouldCallOnFileChange := t.onFileChange != nil && !t.isRunning
		delete(t.debounceFiles, path)
		t.mu.Unlock()

		if shouldCallOnFileChange {
			go t.onFileChange()
		}
	})

	t.mu.Unlock()
}
