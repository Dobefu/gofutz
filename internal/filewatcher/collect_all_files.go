package filewatcher

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// CollectAllFiles collects all files.
func CollectAllFiles() ([]string, error) {
	cwd, err := os.Getwd()

	if err != nil {
		return []string{}, fmt.Errorf(
			"could not get current working directory: %s",
			err.Error(),
		)
	}

	files := []string{}

	err = filepath.Walk(cwd, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		isFile, err := filepath.Match(
			filepath.Join(path, "..", "*.go"),
			path,
		)

		if err != nil {
			return err
		}

		if !isFile {
			return nil
		}

		path = strings.TrimPrefix(path, fmt.Sprintf("%s/", cwd))
		files = append(files, path)

		return nil
	})

	if err != nil {
		return []string{}, fmt.Errorf("could not find any files: %s", err.Error())
	}

	return files, nil
}
