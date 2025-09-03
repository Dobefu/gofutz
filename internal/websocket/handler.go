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

var (
	sharedRunner         *testrunner.TestRunner
	initSharedRunnerOnce sync.Once
)

// Handler defines a websocket handler.
type Handler struct {
	runner *testrunner.TestRunner
	mu     sync.Mutex
	wsChan chan Message
}

// NewHandler creates a new handler.
func NewHandler() (*Handler, error) {
	var err error

	initSharedRunnerOnce.Do(func() {
		files, collectErr := filewatcher.CollectAllFiles()

		if collectErr != nil {
			err = collectErr

			return
		}

		runner, newRunnerErr := testrunner.NewTestRunner(files)

		if newRunnerErr != nil {
			err = newRunnerErr

			return
		}

		sharedRunner = runner
	})

	if err != nil {
		return nil, err
	}

	return &Handler{
		runner: sharedRunner,
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
	files := h.runner.GetFiles()

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
				Files: files,
			},
		})

	case "gofutz:run-all-tests":
		return h.handleRunAllTests(files)

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

func (h *Handler) handleRunAllTests(files map[string]testrunner.File) error {
	for i, file := range files {
		files[i] = testrunner.File{
			Name:            file.Name,
			Functions:       file.Functions,
			Code:            file.Code,
			HighlightedCode: file.HighlightedCode,
			Status:          testrunner.TestStatusRunning,
			Coverage:        -1,
			CoveredLines:    []testrunner.Line{},
		}

		for j, function := range file.Functions {
			files[i].Functions[j] = testrunner.Function{
				Name: function.Name,
				Result: testrunner.TestResult{
					Coverage: -1,
				},
			}
		}
	}

	err := h.SendResponse(Message{
		Method: "gofutz:update",
		Error:  "",
		Params: Params{
			Files: files,
		},
	})

	if err != nil {
		return err
	}

	h.runner.RunAllTests(func(fileToUpdate testrunner.File) error {
		return h.SendResponse(Message{
			Method: "gofutz:update",
			Error:  "",
			Params: Params{
				Files: map[string]testrunner.File{
					fileToUpdate.Name: fileToUpdate,
				},
			},
		})
	})

	return nil
}
