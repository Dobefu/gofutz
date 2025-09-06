package websocket

import (
	"testing"
)

func TestOutputMessage(t *testing.T) {
	t.Parallel()

	message := OutputMessage{
		Method: "gofutz:output",
		Error:  "",
		Params: OutputParams{
			Output: []string{"output"},
		},
	}

	if message.GetMethod() != "gofutz:output" {
		t.Errorf(
			"expected method to be gofutz:output, got: %s",
			message.GetMethod(),
		)
	}

	if message.GetError() != "" {
		t.Errorf("expected error to be empty, got: %s", message.GetError())
	}
}
