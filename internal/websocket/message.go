package websocket

// Message defines a websocket message.
type Message interface {
	GetMethod() string
	GetError() string
}
