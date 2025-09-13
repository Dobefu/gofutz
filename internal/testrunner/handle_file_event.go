package testrunner

import (
	"strings"
	"time"
)

func (t *TestRunner) handleFileEvent(path, operation string) {
	if !strings.HasSuffix(path, ".go") {
		return
	}

	if operation == "CHMOD" || operation == "RENAME" {
		return
	}

	t.mu.Lock()
	timer, hasTimer := t.debounceFiles[path]

	if hasTimer {
		timer.Stop()
	}

	t.debounceFiles[path] = time.AfterFunc(t.debounceDuration, func() {
		t.mu.Lock()

		if t.isRunning {
			delete(t.debounceFiles, path)
			t.mu.Unlock()

			return
		}

		t.mu.Unlock()

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
