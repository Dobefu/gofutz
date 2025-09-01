package websocket

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/Dobefu/gofutz/internal/filewatcher"
	"github.com/Dobefu/gofutz/internal/testrunner"
	"github.com/gorilla/websocket"
)

var runner *testrunner.TestRunner

func init() {
	files, err := filewatcher.CollectAllTestFiles()

	if err != nil {
		slog.Error(err.Error())

		os.Exit(1)
	}

	runner, err = testrunner.NewTestRunner(files)

	if err != nil {
		slog.Error(err.Error())

		os.Exit(1)
	}
}

// Handler defines a websocket handler.
type Handler struct{}

// NewHandler creates a new handler.
func NewHandler() *Handler {
	return &Handler{}
}

// HandleMessage handles a websocket message.
func (h *Handler) HandleMessage(
	ws WsInterface,
	messageType int,
	msg Message,
) error {
	if messageType != websocket.TextMessage {
		return nil
	}

	switch msg.Method {
	case "gofutz:init":
		return h.SendResponse(ws, Message{
			Method: "gofutz:init",
			Error:  "",
			Params: Params{
				Files: runner.GetTests(),
			},
		})

	default:
		return h.SendResponse(ws, Message{
			Method: "error",
			Error:  fmt.Sprintf("Unknown method: %s", msg.Method),
			Params: Params{
				Files: nil,
			},
		})
	}
}

// SendResponse sends a websocket response.
func (h *Handler) SendResponse(ws WsInterface, msg Message) error {
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
