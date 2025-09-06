package testrunner

import (
	"encoding/json"
	"strings"
)

// TestEvent represents a single test event from Go test output in JSON format.
type TestEvent struct {
	Time    string  `json:"time"`
	Action  string  `json:"action"`
	Package string  `json:"package"`
	Test    string  `json:"test"`
	Output  string  `json:"output"`
	Elapsed float64 `json:"elapsed"`
}

func (t *TestRunner) parseFailedTestFilesFromOutput(
	output string,
) map[string]bool {
	failedTestFiles := make(map[string]bool)
	lines := strings.SplitSeq(output, "\n")

	for line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		var event TestEvent
		err := json.Unmarshal([]byte(line), &event)

		if err != nil || event.Action != "fail" {
			continue
		}

		if event.Test == "" || strings.Contains(event.Test, "/") {
			continue
		}

		t.addFailedTestsFromFiles(event, failedTestFiles)
	}

	return failedTestFiles
}

func (t *TestRunner) addFailedTestsFromFiles(
	event TestEvent,
	failedTestFiles map[string]bool,
) {
	for filePath, file := range t.files {
		if !strings.HasPrefix(filePath, event.Package) {
			continue
		}

		for _, function := range file.Functions {
			funcName := strings.TrimPrefix(strings.ToLower(function.Name), "test")
			eventTest := strings.TrimPrefix(strings.ToLower(event.Test), "test")

			if funcName == eventTest {
				failedTestFiles[file.Name] = true

				break
			}
		}
	}
}
