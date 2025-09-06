package testrunner

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// GetModuleName gets the module name of the Go project in the specified path.
func GetModuleName(path string) string {
	file, err := os.Open(filepath.Clean(path))

	if err != nil {
		return ""
	}

	defer func() { _ = file.Close() }()

	scanner := bufio.NewScanner(file)

	if scanner.Scan() {
		line := scanner.Text()
		after, hasAfter := strings.CutPrefix(line, "module ")

		if hasAfter {
			return strings.TrimSpace(after)
		}
	}

	return ""
}
