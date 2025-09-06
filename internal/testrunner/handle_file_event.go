package testrunner

import (
	"strings"
	"time"
)

func (t *TestRunner) handleFileEvent(path, operation string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	timer, hasTimer := t.debounceFiles[path]

	if hasTimer {
		timer.Stop()
	}

	t.debounceFiles[path] = time.AfterFunc(100*time.Millisecond, func() {
		if !strings.HasSuffix(path, ".go") {
			return
		}

		t.processFileEvent(path, operation)

		if t.onFileChange != nil && !t.isRunning {
			go t.onFileChange()
		}

		t.mu.Lock()
		delete(t.debounceFiles, path)
		t.mu.Unlock()
	})
}
