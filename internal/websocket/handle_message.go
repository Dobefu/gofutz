package websocket

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/gorilla/websocket"
)

// HandleMessage handles a websocket message.
func (h *Handler) HandleMessage(
	ws WsInterface,
	messageType int,
	msg Message,
) error {
	if h.runner == nil {
		return fmt.Errorf("test runner is not initialized")
	}

	files := h.runner.GetFiles()

	if messageType != websocket.TextMessage {
		return nil
	}

	if h.wsChan == nil {
		h.wsChan = make(chan Message, 1000)
		go h.handleMessages(ws)
	}

	switch msg.GetMethod() {
	case "gofutz:init":
		if !h.runner.GetHasRunTests() {
			go func() {
				time.Sleep(100 * time.Millisecond)
				err := h.handleRunAllTests()

				if err != nil {
					slog.Error(fmt.Sprintf("Could not run all tests: %s", err.Error()))
				}
			}()
		}

		return h.SendResponse(InitMessage{
			Method: "gofutz:init",
			Error:  "",
			Params: InitParams{
				Files:     files,
				Coverage:  h.runner.GetCoverage(),
				IsRunning: h.runner.GetIsRunning(),
				Output:    h.runner.GetOutput(),
			},
		})

	case "gofutz:run-all-tests":
		return h.handleRunAllTests()

	case "gofutz:stop-tests":
		return h.handleStopTests()

	default:
		return h.SendResponse(UpdateMessage{
			Method: "error",
			Error:  fmt.Sprintf("Unknown method: %s", msg.GetMethod()),
			Params: UpdateParams{
				Files:     nil,
				Coverage:  h.runner.GetCoverage(),
				IsRunning: h.runner.GetIsRunning(),
			},
		})
	}
}

func (h *Handler) handleMessages(ws WsInterface) {
	for msg := range h.wsChan {
		h.mu.Lock()

		json, err := json.Marshal(msg)

		if err != nil {
			h.mu.Unlock()

			continue
		}

		if ws == nil {
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
