package websocket

// InitMessage defines an init message.
type InitMessage struct {
	Method string     `json:"method"`
	Error  string     `json:"error"`
	Params InitParams `json:"params"`
}

// GetMethod returns the method of the message.
func (m InitMessage) GetMethod() string {
	return m.Method
}

// GetError returns the error of the message.
func (m InitMessage) GetError() string {
	return m.Error
}
