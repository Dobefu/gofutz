package ws

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

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

			initMessage := websocket.Message{
				Method: "gofutz:init",
				Error:  "",
				Params: websocket.Params{
					Files: nil,
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
