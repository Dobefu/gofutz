package filewatcher

import (
	"os"
	"strings"
	"testing"
)

func TestAddDirectory(t *testing.T) {
	t.Parallel()

	cwd, err := os.Getwd()

	if err != nil {
		t.Fatalf("expected no error, got: \"%s\"", err.Error())
	}

	tests := []struct {
		name string
	}{
		{
			name: "regular test files",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			fw, err := NewFileWatcher()

			if err != nil {
				t.Fatalf(
					"expected no error creating filewatcher, got: \"%s\"",
					err.Error(),
				)
			}

			defer func() { _ = fw.Close() }()

			err = fw.addDirectory(cwd)

			if err != nil {
				t.Fatalf("expected no error, got: \"%s\"", err.Error())
			}
		})
	}
}

func TestAddDirectoryErr(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		dir      string
		expected string
	}{
		{
			name:     "nonexistent directory",
			dir:      "/bogus",
			expected: "no such file or directory",
		},
		{
			name:     "error reading directory",
			dir:      os.DevNull,
			expected: "not a directory",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			fw, err := NewFileWatcher()

			if err != nil {
				t.Fatalf("expected no error, got: \"%s\"", err.Error())
			}

			defer func() { _ = fw.Close() }()

			err = fw.addDirectory(test.dir)

			if err == nil {
				t.Fatalf("expected error, got nil")
			}

			if !strings.Contains(err.Error(), test.expected) {
				t.Fatalf(
					"expected error to contain \"%s\", got: \"%s\"",
					test.expected,
					err.Error(),
				)
			}
		})
	}
}
