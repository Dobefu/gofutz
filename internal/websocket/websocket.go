// Package websocket provides websocket functionality.
package websocket

import (
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

	return nil
}
