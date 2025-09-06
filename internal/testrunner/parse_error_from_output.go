package testrunner

import (
	"encoding/json"
	"strings"
)

// ParseErrorFromOutput extracts error messages from the output.
func (t *TestRunner) ParseErrorFromOutput(output string) string {
	var errorMessages []string
	lines := strings.SplitSeq(output, "\n")

	for line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		var event TestEvent
		err := json.Unmarshal([]byte(line), &event)

		if err != nil {
			continue
		}

		if strings.Contains(event.Output, "FAIL") ||
			strings.Contains(event.Output, "Error:") ||
			strings.Contains(event.Output, "panic:") ||
			strings.Contains(event.Output, "expected") {
			errorMessages = append(errorMessages, event.Output)
		}
	}

	return strings.Join(errorMessages, "\n")
}
