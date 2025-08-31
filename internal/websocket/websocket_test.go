package websocket

import (
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

const (
	numMsgsBeforeClose = 3
)

type mockWebsocket struct {
	numMsgsRead   int
	isCloseCalled bool
}

func (m *mockWebsocket) SetReadLimit(_ int64) {
	// noop
}

func (m *mockWebsocket) SetReadDeadline(_ time.Time) error {
	return nil
}

func (m *mockWebsocket) SetPongHandler(_ func(string) error) {
	// noop
}

func (m *mockWebsocket) WriteControl(_ int, _ []byte, _ time.Time) error {
	return nil
}

func (m *mockWebsocket) WriteMessage(_ int, _ []byte) error {
	return nil
}

func (m *mockWebsocket) ReadMessage() (int, []byte, error) {
	m.numMsgsRead++

	if m.numMsgsRead > numMsgsBeforeClose {
		return 0, nil, &websocket.CloseError{
			Code: websocket.CloseNormalClosure,
			Text: "",
		}
	}

	return websocket.TextMessage, []byte(`{"type":"test","data":"test"}`), nil
}

func (m *mockWebsocket) Close() error {
	m.isCloseCalled = true

	return nil
}

func TestNewWebsocket(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		ws   *mockWebsocket
	}{
		{
			name: "regular websocket",
			ws: &mockWebsocket{
				numMsgsRead:   0,
				isCloseCalled: false,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			websocket, err := NewWebsocket(test.ws)

			if err != nil {
				t.Fatalf("expected no error, got: %s", err.Error())
			}

			if websocket == nil {
				t.Fatalf("expected websocket, got nil")
			}

			websocket.AddGoroutine()
			go websocket.HandlePing(test.ws)

			err = websocket.HandleMessages(test.ws)

			if err != nil {
				t.Fatalf("expected no error, got: %s", err.Error())
			}

			websocket.FinishGoroutine()
			websocket.Close()

			if !test.ws.isCloseCalled {
				t.Error("expected Close() to have been called on the websocket")
			}
		})
	}
}
