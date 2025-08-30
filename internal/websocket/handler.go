package websocket

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

// Handler defines a websocket handler.
type Handler struct{}

// NewHandler creates a new handler.
func NewHandler() *Handler {
	return &Handler{}
}

// HandleMessage handles a websocket message.
func (h *Handler) HandleMessage(
	ws *websocket.Conn,
	messageType int,
	msg Message,
) error {
	if messageType != websocket.TextMessage {
		return nil
	}

	switch msg.Method {
	default:
		return h.SendResponse(ws, Message{
			Method: "error",
			Params: []any{
				fmt.Sprintf("Unknown method: %s", msg.Method),
			},
		})
	}
}

// SendResponse sends a websocket response.
func (h *Handler) SendResponse(ws *websocket.Conn, msg Message) error {
	json, err := json.Marshal(msg)

	if err != nil {
		return err
	}

	err = ws.WriteMessage(websocket.TextMessage, json)

	if err != nil {
		return err
	}

	return nil
}
