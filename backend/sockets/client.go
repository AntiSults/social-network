package sockets

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	_ "github.com/mattn/go-sqlite3"
)

type ClientList map[*Client]bool

type Client struct {
	connection *websocket.Conn
	manager    *Manager
	egress     chan Event
	nickname   string
}

func NewClient(conn *websocket.Conn, manager *Manager) *Client {
	return &Client{
		connection: conn,
		manager:    manager,
		egress:     make(chan Event),
	}
}

// readMessages reads all possible traffic from WS
func (c *Client) readMessages() {
	var event Event
	defer func() {
		c.manager.removeClient(c)
	}()
	for {
		_, payload, err := c.connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message: %v\n", err)
			}
			break
		}
		err = json.Unmarshal(payload, &event)
		if err != nil {
			log.Println("error unmarshalling JSON:", err)
			continue
		}
		fmt.Printf("Received event: %v\n", event.Type)

		if err := c.manager.routeEvent(event, c); err != nil {
			log.Println("error handling message:", err)
		}
	}
}

func (c *Client) writeMessages() {
	defer func() {
		c.manager.removeClient(c)
	}()
	for message := range c.egress {
		data, err := json.Marshal(message)
		if err != nil {
			log.Println("error marshaling message:", err)
			continue
		}
		if err := c.connection.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Printf("failed to send message: %v", err)
			continue
		}
		log.Println("Message sent:", message.Type)
	}
	if err := c.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
		log.Println("error closing connection:", err)
		return
	}
}
