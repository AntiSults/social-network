package sockets

import "encoding/json"

type Event struct {
	Type         string          `json:"type"`
	Payload      json.RawMessage `json:"payload"`
	SessionToken string          `json:"sessionToken,omitempty"`
}

type EventHandler func(event Event, c *Client) error

// newEvent will be used to create Events for sending to FrontEnd(if any)
func newEvent(t string, p json.RawMessage, s string) *Event {
	return &Event{
		Type:         t,
		Payload:      p,
		SessionToken: s,
	}
}

const (
	EventMessage      = "chat_message"
	EventGroupMessage = "group_chat_message"
	EventUpload       = "initial_upload"
	EventGroupUpload  = "initial_group_upload"
	EventNotify       = "pending_follow_request"
)
