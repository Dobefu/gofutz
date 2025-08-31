// Package ws provides a websocket route handler.
package ws

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/Dobefu/gofutz/internal/websocket"
	gorillawebsocket "github.com/gorilla/websocket"
)

var (
	upgrader = gorillawebsocket.Upgrader{
		HandshakeTimeout:  time.Second * 5,
		ReadBufferSize:    1024,
		WriteBufferSize:   1024,
		WriteBufferPool:   nil,
		Subprotocols:      []string{},
		Error:             nil,
		CheckOrigin:       func(_ *http.Request) bool { return true },
		EnableCompression: true,
	}
)

// Handle handles the route.
func Handle(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		slog.Error(
			fmt.Sprintf("Could not upgrade websocket connection: %s", err.Error()),
		)

		return
	}

	websocketInstance, err := websocket.NewWebsocket(ws)

	if err != nil {
		slog.Error(fmt.Sprintf("Could not create websocket: %s", err.Error()))

		return
	}

	websocketInstance.AddGoroutine()
	go websocketInstance.HandlePing(ws)

	websocketInstance.AddGoroutine()

	err = websocketInstance.HandleMessages(ws)
	websocketInstance.FinishGoroutine()

	if err != nil {
		slog.Error(fmt.Sprintf("Could not handle websocket messages: %s", err.Error()))

		return
	}

	websocketInstance.Close()
}
