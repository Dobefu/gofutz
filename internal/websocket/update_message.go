package websocket

// UpdateMessage defines an update message.
type UpdateMessage struct {
	Method string       `json:"method"`
	Error  string       `json:"error"`
	Params UpdateParams `json:"params"`
}

// GetMethod returns the method of the message.
func (m UpdateMessage) GetMethod() string {
	return m.Method
}

// GetError returns the error of the message.
func (m UpdateMessage) GetError() string {
	return m.Error
}
