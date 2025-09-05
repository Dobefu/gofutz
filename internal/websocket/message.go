package websocket

// Message defines a websocket message.
type Message interface {
	GetMethod() string
	GetError() string
}

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
