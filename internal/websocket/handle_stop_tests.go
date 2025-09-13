package websocket

import (
	"fmt"
)

func (h *Handler) handleStopTests() error {
	if h.runner == nil {
		return fmt.Errorf("test runner is not initialized")
	}

	h.runner.StopTests()
	h.runner.SetIsRunning(false)

	return h.SendResponse(UpdateMessage{
		Method: "gofutz:update",
		Error:  "",
		Params: UpdateParams{
			Files:     nil,
			Coverage:  h.runner.GetCoverage(),
			IsRunning: h.runner.GetIsRunning(),
		},
	})
}
