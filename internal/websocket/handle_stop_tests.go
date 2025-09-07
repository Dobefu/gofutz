package websocket

func (h *Handler) handleStopTests() error {
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
