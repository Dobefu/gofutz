// Package filewatcher provides file watching functionality.
package filewatcher

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// CollectAllTestFiles collects all test files.
func CollectAllTestFiles() ([]string, error) {
	cwd, err := os.Getwd()

	if err != nil {
		return []string{}, fmt.Errorf(
			"could not get current working directory: %s",
			err.Error(),
		)
	}

	testFiles := []string{}

	err = filepath.Walk(cwd, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		isTestFile, err := filepath.Match(
			filepath.Join(path, "..", "*_test.go"),
			path,
		)

		if err != nil {
			return err
		}

		if !isTestFile {
			return nil
		}

		testFiles = append(testFiles, path)

		return nil
	})

	if err != nil {
		return []string{}, fmt.Errorf("could not find any test files: %s", err.Error())
	}

	return testFiles, nil
}
