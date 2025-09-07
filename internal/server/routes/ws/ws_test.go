package ws

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Dobefu/gofutz/internal/websocket"
	gorillawebsocket "github.com/gorilla/websocket"
)

func TestHandle(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
	}{
		{
			name: "regular websocket",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(http.HandlerFunc(Handle))
			defer server.Close()

			ws, _, err := gorillawebsocket.DefaultDialer.Dial(
				fmt.Sprintf("ws%s/ws", server.URL[4:]),
				nil,
			)

			if err != nil {
				t.Fatalf("expected no error, got: %s", err.Error())
			}

			defer func() { _ = ws.Close() }()

			time.Sleep(100 * time.Millisecond)

			initMessage := websocket.UpdateMessage{
				Method: "gofutz:init",
				Error:  "",
				Params: websocket.UpdateParams{
					Files:     nil,
					Coverage:  0,
					IsRunning: false,
				},
			}

			messageBytes, err := json.Marshal(initMessage)

			if err != nil {
				t.Fatalf("could not marshal message: %v", err)
			}

			err = ws.WriteMessage(gorillawebsocket.TextMessage, messageBytes)

			if err != nil {
				t.Fatalf("could not send message: %v", err)
			}
		})
	}
}

func TestHandleErr(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		setup    func(*http.Request)
		expected int
	}{
		{
			name: "invalid headers",
			setup: func(req *http.Request) {
				req.Header.Set("Upgrade", "websocket")
			},
			expected: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest("GET", "/ws", nil)
			test.setup(req)
			w := httptest.NewRecorder()

			Handle(w, req)

			if w.Code != test.expected {
				t.Fatalf(
					"expected status code to be %d, got %d",
					test.expected,
					w.Code,
				)
			}
		})
	}
}
