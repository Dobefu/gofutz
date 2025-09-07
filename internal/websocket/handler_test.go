package websocket

import (
	"testing"
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

func TestSendResponse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		setup func(*Handler)
	}{
		{
			name: "bogus method",
			setup: func(h *Handler) {
				h.wsChan = make(chan Message, 1)
				h.isChannelClosed = false
				h.wsChan <- UpdateMessage{
					Method: "bogus:method",
					Error:  "",
					Params: UpdateParams{
						Files:     nil,
						Coverage:  0,
						IsRunning: false,
					},
				}
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

			test.setup(handler)

			err = handler.SendResponse(UpdateMessage{
				Method: "bogus:method",
				Error:  "",
				Params: UpdateParams{
					Files:     nil,
					Coverage:  0,
					IsRunning: false,
				},
			})

			if err != nil {
				t.Fatalf("expected no error, got: %s", err.Error())
			}
		})
	}
}
