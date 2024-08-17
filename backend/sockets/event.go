package sockets

import "encoding/json"

type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type EventHandler func(event Event, c *Client) error

// newEvent will be used to create Events for sending to FrontEnd(if any)
func newEvent(t string, p json.RawMessage) *Event {
	return &Event{
		Type:    t,
		Payload: p,
	}
}

const (
	EventMessage = "chat_message"
)
