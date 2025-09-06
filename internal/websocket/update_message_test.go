package websocket

import (
	"testing"
)

func TestUpdateMessage(t *testing.T) {
	t.Parallel()

	message := UpdateMessage{
		Method: "gofutz:update",
		Error:  "",
		Params: UpdateParams{
			Files:     nil,
			Coverage:  0,
			IsRunning: false,
		},
	}

	if message.GetMethod() != "gofutz:update" {
		t.Errorf(
			"expected method to be gofutz:update, got: %s",
			message.GetMethod(),
		)
	}

	if message.GetError() != "" {
		t.Errorf("expected error to be empty, got: %s", message.GetError())
	}
}
