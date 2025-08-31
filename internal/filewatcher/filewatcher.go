// Package filewatcher provides file watching functionality.
package filewatcher

import (
	"log/slog"
	"os"

	"github.com/fsnotify/fsnotify"
)

var (
	watcher   *fsnotify.Watcher
	listeners []func(string, string)
)

func init() {
	var err error
	watcher, err = fsnotify.NewWatcher()

	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	go handleFileEvents()

	cwd, err := os.Getwd()

	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	err = watcher.Add(cwd)

	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

func handleFileEvents() {
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			for _, listener := range listeners {
				listener(event.Name, event.Op.String())
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}

			slog.Error(err.Error())
		}
	}
}

// AddListener adds an event listener.
func AddListener(listener func(string, string)) {
	listeners = append(listeners, listener)
}

// Close shuts down the file watcher.
func Close() error {
	if watcher != nil {
		return watcher.Close()
	}

	return nil
}
