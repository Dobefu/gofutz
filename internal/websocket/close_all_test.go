package websocket

import (
	"testing"
)

func TestCloseAll(t *testing.T) {
	t.Parallel()

	handler1 := &Handler{wsChan: make(chan Message, 1)} // nolint:exhaustruct
	handler2 := &Handler{wsChan: make(chan Message, 1)} // nolint:exhaustruct

	handlersMutex.Lock()
	activeHandlers = []*Handler{handler1, handler2}
	handlersMutex.Unlock()

	if len(activeHandlers) != 2 {
		t.Error("there should be 2 active handlers")
	}

	CloseAll()

	handlersMutex.RLock()

	if len(activeHandlers) != 0 {
		t.Error("activeHandlers should be empty")
	}

	handlersMutex.RUnlock()
}
