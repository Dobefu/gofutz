package websocket

import (
	"testing"
)

func TestInitMessage(t *testing.T) {
	t.Parallel()

	message := InitMessage{
		Method: "gofutz:init",
		Error:  "",
		Params: InitParams{
			Files:     nil,
			Coverage:  0,
			IsRunning: false,
			Output:    []string{"output"},
		},
	}

	if message.GetMethod() != "gofutz:init" {
		t.Errorf(
			"expected method to be gofutz:init, got: %s",
			message.GetMethod(),
		)
	}

	if message.GetError() != "" {
		t.Errorf("expected error to be empty, got: %s", message.GetError())
	}
}
