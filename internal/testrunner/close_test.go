package testrunner

import (
	"testing"
)

func TestClose(t *testing.T) {
	t.Parallel()

	runner := &TestRunner{} // nolint:exhaustruct

	runner.Close()
}
