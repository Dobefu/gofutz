package websocket

import (
	"fmt"
	"log/slog"
	"sync"

	"github.com/Dobefu/gofutz/internal/filewatcher"
	"github.com/Dobefu/gofutz/internal/testrunner"
)

var (
	sharedRunner         *testrunner.TestRunner
	initSharedRunnerOnce sync.Once
	initErr              error
	activeHandlers       []*Handler
	handlersMutex        sync.RWMutex
)

// Handler defines a websocket handler.
type Handler struct {
	runner          *testrunner.TestRunner
	mu              sync.Mutex
	wsChan          chan Message
	isChannelClosed bool
}

// NewHandler creates a new handler.
func NewHandler() (*Handler, error) {
	initSharedRunnerOnce.Do(func() {
		files, err := filewatcher.CollectAllFiles()

		if err != nil {
			initErr = fmt.Errorf("could not collect all files: %s", err.Error())

			return
		}

		var sharedWatcher *filewatcher.FileWatcher
		sharedWatcher, err = filewatcher.NewFileWatcher()

		if err != nil {
			initErr = fmt.Errorf("could not create file watcher: %s", err.Error())

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
			initErr = fmt.Errorf("could not initialize runner: %s", err.Error())

			return
		}

		sharedRunner = runner
	})

	if initErr != nil {
		return nil, initErr
	}

	handler := &Handler{
		runner:          sharedRunner,
		mu:              sync.Mutex{},
		wsChan:          nil,
		isChannelClosed: false,
	}

	handlersMutex.Lock()
	activeHandlers = append(activeHandlers, handler)
	handlersMutex.Unlock()

	return handler, nil
}

// SendResponse sends a websocket response.
func (h *Handler) SendResponse(msg Message) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.wsChan == nil || h.isChannelClosed {
		return nil
	}

	h.wsChan <- msg

	return nil
}

func (h *Handler) handleRunAllTests() error {
	if h.runner == nil {
		return fmt.Errorf("test runner is not initialized")
	}

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

// Close closes the handler.
func (h *Handler) Close() {
	h.mu.Lock()

	if h.wsChan != nil && !h.isChannelClosed {
		h.isChannelClosed = true
		close(h.wsChan)
	}

	h.mu.Unlock()

	handlersMutex.Lock()

	for i, handler := range activeHandlers {
		if handler == h {
			activeHandlers = append(activeHandlers[:i], activeHandlers[i+1:]...)

			break
		}
	}

	handlersMutex.Unlock()
}
