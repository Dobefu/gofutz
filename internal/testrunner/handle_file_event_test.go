package testrunner

import (
	"strings"
	"testing"
	"time"
)

func TestHandleFileEvent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		path      string
		operation string
		expected  int
	}{
		{
			name:      "write event",
			path:      "test.go",
			operation: "write",
			expected:  1,
		},
		{
			name:      "non-go file event",
			path:      "test",
			operation: "write",
			expected:  1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			timeoutFired := make(chan struct{})
			runner := &TestRunner{ // nolint:exhaustruct
				debounceFiles:    make(map[string]*time.Timer),
				debounceDuration: 1 * time.Millisecond,
				onFileChange:     func() { close(timeoutFired) },
			}

			runner.handleFileEvent(test.path, test.operation)
			runner.handleFileEvent(test.path, test.operation)

			if len(runner.debounceFiles) != test.expected {
				t.Fatalf(
					"expected %d debounce file(s), got: %d",
					test.expected,
					len(runner.debounceFiles),
				)
			}

			if !strings.HasSuffix(test.path, ".go") {
				return
			}

			select {
			case <-timeoutFired:
			case <-time.After(100 * time.Millisecond):
				t.Fatal("timeout did not fire within the expected timeframe")
			}

			time.Sleep(1 * time.Millisecond)

			if len(runner.debounceFiles) != 0 {
				t.Fatalf(
					"expected 0 debounce files after timeout, got: %d",
					len(runner.debounceFiles),
				)
			}
		})
	}
}
