package filewatcher

import (
	"os"
	"slices"
	"testing"
)

func TestFilewatcher(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		expected []string
	}{
		{
			name: "regular test files",
			expected: []string{
				"collect_all_files.go",
				"filewatcher.go",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			files, err := CollectAllFiles()

			if err != nil {
				t.Fatalf("expected no error, got: %s", err.Error())
			}

			if len(files) == 0 {
				t.Fatalf("expected files, got none")
			}

			for _, file := range files {
				if !slices.Contains(test.expected, file) {
					t.Errorf("expected file %s, got %s", file, test.expected)
				}
			}
		})
	}
}

func TestAddListener(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		expected int
	}{
		{
			name:     "regular test files",
			expected: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			fw, err := NewFileWatcher()

			if err != nil {
				t.Fatalf("expected no error creating filewatcher, got: %s", err.Error())
			}

			defer func() { _ = fw.Close() }()

			fw.AddListener(func(_ string, _ string) {})

			if fw.GetListenerCount() != test.expected {
				t.Fatalf(
					"expected %d listener(s), got %d",
					test.expected,
					fw.GetListenerCount(),
				)
			}
		})
	}
}

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

func TestMultipleListeners(t *testing.T) {
	t.Parallel()

	fw, err := NewFileWatcher()

	if err != nil {
		t.Fatalf("expected no error creating filewatcher, got: %s", err.Error())
	}

	defer func() { _ = fw.Close() }()

	fw.AddListener(func(_ string, _ string) {})
	fw.AddListener(func(_ string, _ string) {})
	fw.AddListener(func(_ string, _ string) {})

	if fw.GetListenerCount() != 3 {
		t.Fatalf("expected 3 listeners, got %d", fw.GetListenerCount())
	}
}

func TestResetListeners(t *testing.T) {
	t.Parallel()

	fw, err := NewFileWatcher()

	if err != nil {
		t.Fatalf("expected no error creating filewatcher, got: %s", err.Error())
	}

	defer func() { _ = fw.Close() }()

	fw.AddListener(func(_ string, _ string) {})
	fw.AddListener(func(_ string, _ string) {})

	if fw.GetListenerCount() != 2 {
		t.Fatalf("expected 2 listeners, got %d", fw.GetListenerCount())
	}

	fw.ResetListeners()

	if fw.GetListenerCount() != 0 {
		t.Fatalf("expected 0 listeners after reset, got %d", fw.GetListenerCount())
	}
}
