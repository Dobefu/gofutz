package testrunner

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// GetTestsFromFile gets all test functions from a file.
func GetTestsFromFile(file string) ([]string, error) {
	tests := []string{}

	if file == "" {
		return []string{}, fmt.Errorf("file is empty")
	}

	fileContent, err := os.ReadFile(filepath.Clean(file))

	if err != nil {
		return []string{}, err
	}

	lines := strings.SplitSeq(string(fileContent), "\n")

	for line := range lines {
		line = strings.TrimSpace(line)

		if !strings.HasPrefix(line, "func Test") {
			continue
		}

		funcDefinition := strings.TrimPrefix(line, "func ")
		argsIdx := strings.Index(funcDefinition, "(")

		if argsIdx == -1 {
			continue
		}

		funcName := funcDefinition[:argsIdx]

		genericIdx := strings.Index(funcName, "[")

		if genericIdx != -1 {
			funcName = funcName[:genericIdx]
		}

		tests = append(tests, funcName)
	}

	return tests, nil
}
