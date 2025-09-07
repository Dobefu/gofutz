package filewatcher

import (
	"os"
	"path/filepath"
	"strings"
)

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
