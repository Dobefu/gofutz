package websocket

// Message defines a websocket message.
type Message struct {
	Method string `json:"method"`
	Params []any  `json:"params"`
}
