// Package filewatcher provides file watching functionality.
package filewatcher

import (
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
)

// FileWatcher defines a file watcher.
type FileWatcher struct {
	watcher   *fsnotify.Watcher
	listeners []func(string, string)
	mu        sync.RWMutex
	stopCh    chan struct{}
}

// NewFileWatcher creates a new FileWatcher instance.
func NewFileWatcher() (*FileWatcher, error) {
	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		return nil, err
	}

	fw := &FileWatcher{
		watcher:   watcher,
		listeners: make([]func(string, string), 0),
		mu:        sync.RWMutex{},
		stopCh:    make(chan struct{}),
	}

	go fw.handleFileEvents()
	cwd, err := os.Getwd()

	if err != nil {
		_ = fw.Close()

		return nil, err
	}

	err = fw.addDirectory(cwd)

	if err != nil {
		_ = fw.Close()

		return nil, err
	}

	return fw, nil
}

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

// AddListener adds an event listener.
func (fw *FileWatcher) AddListener(listener func(string, string)) {
	fw.mu.Lock()
	fw.listeners = append(fw.listeners, listener)
	fw.mu.Unlock()
}

// GetListenerCount returns the current number of listeners.
func (fw *FileWatcher) GetListenerCount() int {
	fw.mu.RLock()
	count := len(fw.listeners)
	fw.mu.RUnlock()

	return count
}

// ResetListeners resets all listeners.
func (fw *FileWatcher) ResetListeners() {
	fw.mu.Lock()
	fw.listeners = nil
	fw.mu.Unlock()
}

// addDirectory adds all directories in the given directory to the watcher.
func (fw *FileWatcher) addDirectory(dir string) error {
	err := fw.watcher.Add(dir)

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
		err = fw.addDirectory(subDir)

		if err != nil {
			return err
		}
	}

	return nil
}

// Close shuts down the file watcher.
func (fw *FileWatcher) Close() error {
	close(fw.stopCh)

	if fw.watcher != nil {
		return fw.watcher.Close()
	}

	return nil
}
