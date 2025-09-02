package websocket

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"

	"github.com/Dobefu/gofutz/internal/filewatcher"
	"github.com/Dobefu/gofutz/internal/testrunner"
	"github.com/gorilla/websocket"
)

// Handler defines a websocket handler.
type Handler struct {
	runner *testrunner.TestRunner
	mu     sync.Mutex
	wsChan chan Message
}

// NewHandler creates a new handler.
func NewHandler() (*Handler, error) {
	files, err := filewatcher.CollectAllFiles()

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
		wsChan: nil,
	}, nil
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

	if h.wsChan == nil {
		h.wsChan = make(chan Message, 100)
		go h.handleMessages(ws)
	}

	switch msg.Method {
	case "gofutz:init":
		return h.SendResponse(Message{
			Method: "gofutz:init",
			Error:  "",
			Params: Params{
				Files: h.runner.GetFiles(),
			},
		})

	case "gofutz:run-all-tests":
		h.runner.RunAllTests(func(test testrunner.Test) error {
			return h.SendResponse(Message{
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
		return h.SendResponse(Message{
			Method: "error",
			Error:  fmt.Sprintf("Unknown method: %s", msg.Method),
			Params: Params{
				Files: nil,
			},
		})
	}
}

// SendResponse sends a websocket response.
func (h *Handler) SendResponse(msg Message) error {
	h.wsChan <- msg

	return nil
}

func (h *Handler) handleMessages(ws WsInterface) {
	for msg := range h.wsChan {
		h.mu.Lock()

		json, err := json.Marshal(msg)

		if err != nil {
			h.mu.Unlock()

			continue
		}

		err = ws.WriteMessage(websocket.TextMessage, json)
		h.mu.Unlock()

		if err != nil {
			slog.Error(fmt.Sprintf("Could not send message: %s", err.Error()))

			return
		}
	}
}
