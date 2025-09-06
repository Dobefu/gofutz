package testrunner

import (
	"testing"

	"github.com/Dobefu/gofutz/internal/filewatcher"
)

func TestNewTestRunner(t *testing.T) {
	t.Parallel()

	fw, err := filewatcher.NewFileWatcher()

	if err != nil {
		t.Fatalf("expected no error creating filewatcher, got: %s", err.Error())
	}

	defer func() { _ = fw.Close() }()

	runner, err := NewTestRunner([]string{}, fw, func() {})

	if err != nil {
		t.Fatalf("expected no error, got: %s", err.Error())
	}

	if runner == nil {
		t.Fatalf("expected runner to be set, got: nil")
	}
}
