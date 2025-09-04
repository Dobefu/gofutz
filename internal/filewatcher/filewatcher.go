// Package filewatcher provides file watching functionality.
package filewatcher

import (
	"log/slog"
	"os"
	"path/filepath"
	"strings"

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

	err = addDirectory(cwd)

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

// addDirectory adds all directories in the given directory to the watcher.
func addDirectory(dir string) error {
	err := watcher.Add(dir)

	if err != nil {
		return err
	}

	entries, err := os.ReadDir(dir)

	if err != nil {
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		subDir := filepath.Join(dir, entry.Name())
		err = addDirectory(subDir)

		if err != nil {
			return err
		}
	}

	return nil
}

// Close shuts down the file watcher.
func Close() error {
	if watcher != nil {
		return watcher.Close()
	}

	return nil
}
