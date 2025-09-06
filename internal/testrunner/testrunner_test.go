package testrunner

import (
	"testing"
)

func TestNewTestRunner(t *testing.T) {
	t.Parallel()

	runner, err := NewTestRunner([]string{}, func() {})

	if err != nil {
		t.Fatalf("expected no error, got: %s", err.Error())
	}

	if runner == nil {
		t.Fatalf("expected runner to be set, got: nil")
	}
}
