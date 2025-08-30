// Package websocket provides websocket functionality.
package websocket

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/gorilla/websocket"
)

// NewWebsocket creates a new websocket.
func NewWebsocket(ws *websocket.Conn) error {
	defer func() { _ = ws.Close() }()

	ws.SetReadLimit(512)

	err := ws.SetReadDeadline(time.Now().Add(time.Second * 5))

	if err != nil {
		return err
	}

	ws.SetPongHandler(func(string) error {
		err = ws.SetReadDeadline(time.Now().Add(time.Second * 5))

		if err != nil {
			return err
		}

		return nil
	})

	err = ws.WriteMessage(websocket.PongMessage, nil)

	if err != nil {
		return err
	}

	handler := NewHandler()

	for {
		messageType, message, err := ws.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(
				err,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
			) {
				slog.Error(fmt.Sprintf("Could not read message: %s", err.Error()))
			}

			break
		}

		msg := &Message{} // nolint:exhaustruct
		err = json.Unmarshal(message, &msg)

		if err != nil {
			slog.Error(fmt.Sprintf("Could not unmarshal message: %s", err.Error()))

			break
		}

		err = handler.HandleMessage(ws, messageType, *msg)

		if err != nil {
			slog.Error(fmt.Sprintf("Could not handle message: %s", err.Error()))
		}
	}

	return nil
}
