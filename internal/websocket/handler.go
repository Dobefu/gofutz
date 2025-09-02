package websocket

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/Dobefu/gofutz/internal/filewatcher"
	"github.com/Dobefu/gofutz/internal/testrunner"
	"github.com/gorilla/websocket"
)

// Handler defines a websocket handler.
type Handler struct {
	runner *testrunner.TestRunner
	mu     sync.Mutex
}

// NewHandler creates a new handler.
func NewHandler() (*Handler, error) {
	files, err := filewatcher.CollectAllTestFiles()

	if err != nil {
		return nil, err
	}

	runner, err := testrunner.NewTestRunner(files)

	if err != nil {
		return nil, err
	}

	return &Handler{
		runner: runner,
		mu:     sync.Mutex{},
	}, nil
}

// HandleMessage handles a websocket message.
func (h *Handler) HandleMessage(
	ws WsInterface,
	messageType int,
	msg Message,
) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if messageType != websocket.TextMessage {
		return nil
	}

	switch msg.Method {
	case "gofutz:init":
		return h.SendResponse(ws, Message{
			Method: "gofutz:init",
			Error:  "",
			Params: Params{
				Files: h.runner.GetTests(),
			},
		})

	case "gofutz:run-all-tests":
		h.runner.RunAllTests(func(test testrunner.Test) error {
			h.mu.Lock()
			defer h.mu.Unlock()

			return h.SendResponse(ws, Message{
				Method: "gofutz:update",
				Error:  "",
				Params: Params{
					Files: map[string]testrunner.File{
						test.Name: {
							Name:            test.Name,
							Tests:           []testrunner.Test{test},
							Code:            "",
							HighlightedCode: "",
						},
					},
				},
			})
		})

		return nil

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
