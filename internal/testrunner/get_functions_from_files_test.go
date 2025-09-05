package testrunner

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestGetFunctionsFromFiles(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
	}{
		{
			name: "success",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			cwd, err := os.Getwd()

			if err != nil {
				t.Fatalf("expected no error, got: %s", err.Error())
			}

			functions, err := GetFunctionsFromFiles(
				[]string{filepath.Join(cwd, "get_functions_from_files.go")},
			)

			if err != nil {
				t.Fatalf("expected no error, got: %s", err.Error())
			}

			if len(functions) == 0 {
				t.Fatalf("expected functions, got none")
			}
		})
	}
}

func TestGetFunctionsFromFilesErr(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		files    []string
		expected string
	}{
		{
			name:  "nonexistent file",
			files: []string{os.DevNull},
			expected: fmt.Sprintf(
				"%s:1:1: expected 'package', found 'EOF'",
				os.DevNull,
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			_, err := GetFunctionsFromFiles(test.files)

			if err == nil {
				t.Fatalf("expected error, got nil")
			}

			if err.Error() != test.expected {
				t.Fatalf(
					"expected error to be \"%s\", got: \"%s\"",
					test.expected,
					err.Error(),
				)
			}
		})
	}
}
