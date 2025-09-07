package websocket

// CloseAll closes all handlers and shared resources.
func CloseAll() {
	handlersMutex.Lock()
	defer handlersMutex.Unlock()

	for _, handler := range activeHandlers {
		handler.mu.Lock()

		if handler.wsChan != nil {
			handler.isChannelClosed = true
			close(handler.wsChan)
		}

		handler.mu.Unlock()
	}

	activeHandlers = nil

	if sharedRunner != nil {
		sharedRunner.Close()
		sharedRunner = nil
	}
}
