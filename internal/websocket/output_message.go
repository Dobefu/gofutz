package websocket

// OutputMessage defines an output message.
type OutputMessage struct {
	Method string       `json:"method"`
	Error  string       `json:"error"`
	Params OutputParams `json:"params"`
}

// GetMethod returns the method of the message.
func (m OutputMessage) GetMethod() string {
	return m.Method
}

// GetError returns the error of the message.
func (m OutputMessage) GetError() string {
	return m.Error
}
