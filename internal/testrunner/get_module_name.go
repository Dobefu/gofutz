package testrunner

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// GetModuleName gets the module name of the Go project in the current directory.
func GetModuleName() string {
	cwd, err := os.Getwd()

	if err != nil {
		return ""
	}

	goModPath := filepath.Join(cwd, "go.mod")
	file, err := os.Open(filepath.Clean(goModPath))

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
