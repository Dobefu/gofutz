package testrunner

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
)

func (t *TestRunner) processFileEvent(path, operation string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if strings.HasSuffix(path, "_test.go") {
		return
	}

	cwd, err := os.Getwd()

	if err != nil {
		slog.Error(err.Error())

		return
	}

	moduleName := GetModuleName()
	modulePath := fmt.Sprintf("%s%s", moduleName, strings.TrimPrefix(path, cwd))

	switch operation {
	case "CREATE":
		delete(t.files, modulePath)

		fallthrough

	case "WRITE", "MODIFY", "RENAME":
		functions, code, err := GetFunctionsFromFile(path)

		if err != nil {
			slog.Error(err.Error())

			return
		}

		if len(functions) == 0 {
			delete(t.files, modulePath)

			return
		}

		t.files[modulePath] = File{
			Name:            modulePath,
			Functions:       functions,
			Code:            code,
			HighlightedCode: HighlightCode("go", string(code)),
			Status:          TestStatusPending,
			Coverage:        -1,
			CoveredLines:    []Line{},
		}

	case "REMOVE":
		delete(t.files, modulePath)
	}
}
