package filewatcher

import (
	"log/slog"
)

func (fw *FileWatcher) handleFileEvents() {
	for {
		select {
		case event, ok := <-fw.watcher.Events:
			if !ok {
				return
			}

			fw.mu.RLock()

			for _, listener := range fw.listeners {
				listener(event.Name, event.Op.String())
			}

			fw.mu.RUnlock()

		case err, ok := <-fw.watcher.Errors:
			if !ok {
				return
			}

			slog.Error(err.Error())

		case <-fw.stopCh:
			return
		}
	}
}
