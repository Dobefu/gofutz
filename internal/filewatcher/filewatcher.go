// Package filewatcher provides file watching functionality.
package filewatcher

import (
	"os"
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

// Close shuts down the file watcher.
func (fw *FileWatcher) Close() error {
	close(fw.stopCh)

	if fw.watcher != nil {
		return fw.watcher.Close()
	}

	return nil
}
