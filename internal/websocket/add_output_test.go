package websocket

import (
	"testing"

	"github.com/Dobefu/gofutz/internal/testrunner"
)

func TestAddOutput(t *testing.T) {
	t.Parallel()

	handler := &Handler{ // nolint:exhaustruct
		runner: &testrunner.TestRunner{}, // nolint:exhaustruct
	}

	err := handler.AddOutput("output")

	if err != nil {
		t.Fatalf("expected no error, got: \"%s\"", err.Error())
	}

	if len(handler.runner.GetOutput()) != 1 {
		t.Fatalf(
			"expected 1 output line, got: %d",
			len(handler.runner.GetOutput()),
		)
	}

	err = handler.AddOutput("")

	if err != nil {
		t.Fatalf("expected no error, got: \"%s\"", err.Error())
	}

	if len(handler.runner.GetOutput()) != 1 {
		t.Fatalf(
			"expected 1 output line, got: %d",
			len(handler.runner.GetOutput()),
		)
	}
}
