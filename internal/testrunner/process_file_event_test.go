package testrunner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestProcessFileEvent(t *testing.T) {
	t.Parallel()

	cwd, err := os.Getwd()

	if err != nil {
		t.Fatalf("expected no error, got: %s", err.Error())
	}

	tests := []struct {
		name      string
		path      string
		operation string
		expected  int
	}{
		{
			name:      "create event",
			path:      filepath.Join(cwd, "process_file_event.go"),
			operation: "CREATE",
			expected:  1,
		},
		{
			name:      "test file",
			path:      filepath.Join(cwd, "process_file_event_test.go"),
			operation: "WRITE",
			expected:  0,
		},
		{
			name:      "file without functions",
			path:      filepath.Join(cwd, "file.go"),
			operation: "RENAME",
			expected:  0,
		},
		{
			name:      "remove event",
			path:      filepath.Join(cwd, "process_file_event.go"),
			operation: "REMOVE",
			expected:  0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			runner := &TestRunner{ // nolint:exhaustruct
				files: map[string]File{},
			}

			runner.processFileEvent(test.path, test.operation)

			if len(runner.files) != test.expected {
				t.Fatalf(
					"expected %d file(s), got: %d",
					test.expected,
					len(runner.files),
				)
			}
		})
	}
}
