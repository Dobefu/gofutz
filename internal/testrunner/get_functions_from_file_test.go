package testrunner

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestGetFunctionsFromFile(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		expected int
	}{
		{
			name:     "success",
			expected: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			cwd, err := os.Getwd()

			if err != nil {
				t.Fatalf("expected no error, got: %s", err.Error())
			}

			functions, _, err := GetFunctionsFromFile(
				filepath.Join(cwd, "get_functions_from_file.go"),
			)

			if err != nil {
				t.Fatalf("expected no error, got: %s", err.Error())
			}

			if len(functions) != test.expected {
				t.Fatalf(
					"expected %d function(s), got: %d",
					test.expected,
					len(functions),
				)
			}
		})
	}
}

func TestGetFunctionsFromFileErr(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		file     string
		expected string
	}{
		{
			name: "nonexistent file",
			file: os.DevNull,
			expected: fmt.Sprintf(
				"%s:1:1: expected 'package', found 'EOF'",
				os.DevNull,
			),
		},
		{
			name:     "empty filename",
			file:     "",
			expected: "file is empty",
		},
		{
			name:     "error reading file",
			file:     "/bogus",
			expected: "open /bogus: no such file or directory",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
		})

		_, _, err := GetFunctionsFromFile(
			test.file,
		)

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
	}
}
