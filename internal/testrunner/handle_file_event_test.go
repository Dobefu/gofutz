package testrunner

import (
	"strings"
	"sync"
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
			operation: "WRITE",
			expected:  1,
		},
		{
			name:      "non-go file event",
			path:      "test",
			operation: "WRITE",
			expected:  1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			timeoutFired := make(chan struct{})
			debounceCount := 0
			debounceMutex := &sync.Mutex{}

			runner := &TestRunner{ // nolint:exhaustruct
				debounceFiles:    make(map[string]*time.Timer),
				debounceDuration: 1 * time.Millisecond,
				onFileChange: func() {
					debounceMutex.Lock()
					debounceCount++
					debounceMutex.Unlock()

					close(timeoutFired)
				},
			}

			runner.handleFileEvent(test.path, test.operation)
			runner.handleFileEvent(test.path, test.operation)

			time.Sleep(2 * time.Millisecond)

			if !strings.HasSuffix(test.path, ".go") {
				return
			}

			select {
			case <-timeoutFired:
			case <-time.After(100 * time.Millisecond):
				t.Fatal("timeout did not fire within the expected timeframe")
			}

			time.Sleep(2 * time.Millisecond)

			debounceMutex.Lock()
			actualCount := debounceCount
			debounceMutex.Unlock()

			if actualCount != test.expected {
				t.Fatalf(
					"expected %d debounce callbacks, got: %d",
					test.expected,
					actualCount,
				)
			}
		})
	}
}
