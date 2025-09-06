package websocket

import (
	"testing"

	"github.com/gorilla/websocket"
)

func TestNewHandler(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
	}{
		{
			name: "basic handler",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			handler, err := NewHandler()

			if err != nil {
				t.Fatalf("expected no error, got: %s", err.Error())
			}

			if handler == nil {
				t.Fatalf("expected handler, got nil")
			}
		})
	}
}

func TestHandleMessage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		message Message
	}{
		{
			name: "init message",
			message: UpdateMessage{
				Method: "gofutz:init",
				Error:  "",
				Params: UpdateParams{
					Files:     nil,
					Coverage:  0,
					IsRunning: false,
				},
			},
		},
		{
			name: "run all tests message",
			message: UpdateMessage{
				Method: "gofutz:run-all-tests",
				Error:  "",
				Params: UpdateParams{
					Files:     nil,
					Coverage:  0,
					IsRunning: false,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			handler, err := NewHandler()

			if err != nil {
				t.Fatalf("expected no error, got: %s", err.Error())
			}

			if handler == nil {
				t.Fatalf("expected handler, got nil")
			}

			err = handler.HandleMessage(nil, websocket.TextMessage, test.message)

			if err != nil {
				t.Fatalf("expected no error, got: %s", err.Error())
			}
		})
	}
}
