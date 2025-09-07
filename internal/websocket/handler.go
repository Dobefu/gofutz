package websocket

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/Dobefu/gofutz/internal/filewatcher"
	"github.com/Dobefu/gofutz/internal/testrunner"
	"github.com/gorilla/websocket"
)

var (
	sharedRunner         *testrunner.TestRunner
	initSharedRunnerOnce sync.Once
	activeHandlers       []*Handler
	handlersMutex        sync.RWMutex
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

		var sharedWatcher *filewatcher.FileWatcher
		sharedWatcher, err = filewatcher.NewFileWatcher()

		if err != nil {
			return
		}

		runner, err := testrunner.NewTestRunner(files, sharedWatcher, func() {
			handlersMutex.RLock()

			for _, handler := range activeHandlers {
				go func(h *Handler) {
					err := h.handleRunAllTests()

					if err != nil {
						slog.Error(
							fmt.Sprintf(
								"Could not run tests on file change: %s",
								err.Error(),
							),
						)
					}
				}(handler)
			}

			handlersMutex.RUnlock()
		})

		if err != nil {
			_ = sharedWatcher.Close()

			return
		}

		sharedRunner = runner
	})

	if err != nil {
		return nil, err
	}

	handler := &Handler{
		runner: sharedRunner,
		mu:     sync.Mutex{},
		wsChan: nil,
	}

	handlersMutex.Lock()
	activeHandlers = append(activeHandlers, handler)
	handlersMutex.Unlock()

	return handler, nil
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

// SendResponse sends a websocket response.
func (h *Handler) SendResponse(msg Message) error {
	if h.wsChan == nil {
		return nil
	}

	select {
	case h.wsChan <- msg:
		return nil
	default:
		return nil
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

func (h *Handler) handleRunAllTests() error {
	h.runner.SetHasRunTests(true)
	h.runner.SetIsRunning(true)

	files := h.runner.GetFiles()
	newFiles := make(map[string]testrunner.File)

	for name, file := range files {
		var status testrunner.TestStatus = testrunner.TestStatusRunning

		if len(file.Functions) == 0 {
			status = file.Status
		}

		newFile := testrunner.File{
			Name:            file.Name,
			Functions:       make([]testrunner.Function, len(file.Functions)),
			Code:            file.Code,
			HighlightedCode: file.HighlightedCode,
			Status:          status,
			Coverage:        -1,
			CoveredLines:    []testrunner.Line{},
		}

		for j, function := range file.Functions {
			newFile.Functions[j] = testrunner.Function{
				Name: function.Name,
				Result: testrunner.TestResult{
					Coverage: -1,
				},
			}
		}

		newFiles[name] = newFile
	}

	err := h.SendResponse(UpdateMessage{
		Method: "gofutz:update",
		Error:  "",
		Params: UpdateParams{
			Files:     newFiles,
			Coverage:  -1,
			IsRunning: h.runner.GetIsRunning(),
		},
	})

	if err != nil {
		return err
	}

	h.runner.RunAllTests(func(fileToUpdate testrunner.File) error {
		return h.SendResponse(UpdateMessage{
			Method: "gofutz:update",
			Error:  "",
			Params: UpdateParams{
				Files: map[string]testrunner.File{
					fileToUpdate.Name: fileToUpdate,
				},
				Coverage:  h.runner.GetCoverage(),
				IsRunning: h.runner.GetIsRunning(),
			},
		})
	}, func(output string) error {
		return h.AddOutput(output)
	}, func() {
		h.runner.SetIsRunning(false)
	})

	return nil
}

func (h *Handler) handleStopTests() error {
	h.runner.StopTests()
	h.runner.SetIsRunning(false)

	return h.SendResponse(UpdateMessage{
		Method: "gofutz:update",
		Error:  "",
		Params: UpdateParams{
			Files:     nil,
			Coverage:  h.runner.GetCoverage(),
			IsRunning: h.runner.GetIsRunning(),
		},
	})
}

// Close closes the handler.
func (h *Handler) Close() {
	handlersMutex.Lock()

	for i, handler := range activeHandlers {
		if handler == h {
			activeHandlers = append(activeHandlers[:i], activeHandlers[i+1:]...)

			break
		}
	}

	handlersMutex.Unlock()
}
