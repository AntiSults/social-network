package sockets

import "encoding/json"

type Event struct {
	Type         string          `json:"type"`
	Payload      json.RawMessage `json:"payload"`
	SessionToken string          `json:"sessionToken,omitempty"`
}

type EventHandler func(event Event, c *Client) error

// newEvent for creating new event to send message to Front via websocket
func newEvent(t string, p json.RawMessage) *Event {
	return &Event{
		Type:    t,
		Payload: p,
		// SessionToken: s,
	}
}

const (
	EventMessage = "newPM"
)
