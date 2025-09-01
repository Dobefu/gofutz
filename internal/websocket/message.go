package websocket

// Message defines a websocket message.
type Message struct {
	Method string `json:"method"`
	Error  string `json:"error"`
	Params Params `json:"params"`
}
