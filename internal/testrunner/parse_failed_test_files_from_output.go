package testrunner

import (
	"encoding/json"
	"fmt"
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

		packageParts := strings.Split(event.Package, "/")

		if len(packageParts) < 2 {
			continue
		}

		packageName := packageParts[len(packageParts)-1]
		sourceFileName := fmt.Sprintf("%s.go", packageName)

		for filePath := range t.files {
			if strings.HasSuffix(filePath, fmt.Sprintf("/%s", sourceFileName)) {
				failedTestFiles[filePath] = true

				break
			}
		}
	}

	return failedTestFiles
}
