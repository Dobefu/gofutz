package websocket

// CloseAll closes all handlers and shared resources.
func CloseAll() {
	handlersMutex.Lock()
	defer handlersMutex.Unlock()

	for _, handler := range activeHandlers {
		if handler.wsChan != nil {
			close(handler.wsChan)
		}
	}

	activeHandlers = nil

	if sharedRunner != nil {
		sharedRunner.Close()
		sharedRunner = nil
	}
}
