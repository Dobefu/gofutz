package websocket

import (
	"strings"
)

// AddOutput adds to the output buffer and sends it.
func (h *Handler) AddOutput(output string) error {
	outputLines := strings.Split(output, "\n")

	if len(outputLines) > 0 && outputLines[len(outputLines)-1] == "" {
		outputLines = outputLines[:len(outputLines)-1]
	}

	newOutput, err := h.runner.AddOutput(outputLines)

	if err != nil {
		return err
	}

	return h.SendResponse(OutputMessage{
		Method: "gofutz:output",
		Error:  "",
		Params: OutputParams{
			Output: newOutput,
		},
	})
}
