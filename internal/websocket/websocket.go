// Package websocket provides websocket functionality.
package websocket

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// WsInterface defines a websocket interface.
type WsInterface interface {
	SetReadLimit(limit int64)
	SetReadDeadline(deadline time.Time) error
	SetPongHandler(handler func(string) error)
	WriteControl(messageType int, data []byte, deadline time.Time) error
	WriteMessage(messageType int, data []byte) error
	ReadMessage() (messageType int, p []byte, err error)
	Close() error
}

// Websocket defines a websocket.
type Websocket struct {
	handler *Handler
	ws      WsInterface
	close   chan struct{}
	wg      *sync.WaitGroup
}

// NewWebsocket creates a new websocket.
func NewWebsocket(ws WsInterface) (*Websocket, error) {
	handler, err := NewHandler()

	if err != nil {
		return nil, err
	}

	w := &Websocket{
		handler: handler,
		ws:      ws,
		close:   make(chan struct{}),
		wg:      &sync.WaitGroup{},
	}

	ws.SetReadLimit(512)

	err = ws.SetReadDeadline(time.Now().Add(time.Second * 15))

	if err != nil {
		return nil, err
	}

	ws.SetPongHandler(func(string) error {
		err = ws.SetReadDeadline(time.Now().Add(time.Second * 15))

		if err != nil {
			return err
		}

		return nil
	})

	return w, nil
}

// AddGoroutine adds a goroutine to the wait group.
func (w *Websocket) AddGoroutine() {
	w.wg.Add(1)
}

// FinishGoroutine finishes a goroutine after the wait group is done.
func (w *Websocket) FinishGoroutine() {
	w.wg.Done()
}

// HandlePing handles sending pings to the client.
func (w *Websocket) HandlePing(ws WsInterface) {
	defer w.FinishGoroutine()

	isInitialPing := true
	ticker := time.NewTicker(time.Millisecond * 50)
	defer ticker.Stop()

	for range ticker.C {
		if isInitialPing {
			ticker.Reset(time.Second * 5)
			isInitialPing = false
		}

		select {
		case _, ok := <-w.close:
			if !ok {
				return
			}

		default:
		}

		err := ws.WriteControl(
			websocket.PingMessage,
			[]byte{},
			time.Now().Add(time.Second*10),
		)

		if err != nil {
			slog.Error(fmt.Sprintf("Could not send ping: %s", err.Error()))
		}
	}
}

// HandleMessages handles websocket messages.
func (w *Websocket) HandleMessages(ws WsInterface) error {
	for {
		select {
		case _, ok := <-w.close:
			if !ok {
				return nil
			}

		default:
		}

		messageType, message, err := ws.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(
				err,
				websocket.CloseNormalClosure,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
			) {
				slog.Error(fmt.Sprintf("Could not read message: %s", err.Error()))
			}

			break
		}

		msg := &UpdateMessage{} // nolint:exhaustruct
		err = json.Unmarshal(message, &msg)

		if err != nil {
			slog.Error(fmt.Sprintf("Could not unmarshal message: %s", err.Error()))

			break
		}

		err = w.handler.HandleMessage(ws, messageType, *msg)

		if err != nil {
			slog.Error(fmt.Sprintf("Could not handle message: %s", err.Error()))

			break
		}
	}

	return nil
}

// Close sends a signal to close the websocket.
func (w *Websocket) Close() {
	close(w.close)
	w.wg.Wait()

	_ = w.ws.Close()
}
