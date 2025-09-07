package websocket

import (
	"sync"
	"testing"

	"github.com/Dobefu/gofutz/internal/testrunner"
)

func TestHandleStopTests(t *testing.T) {
	t.Parallel()

	handler := &Handler{
		runner:          &testrunner.TestRunner{},
		mu:              sync.Mutex{},
		wsChan:          make(chan Message, 100),
		isChannelClosed: false,
	}

	err := handler.handleStopTests()

	if err != nil {
		t.Fatalf("expected no error, got: %s", err.Error())
	}
}
