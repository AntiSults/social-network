package sockets

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"social-network/structs"
	"sync"
)

var WebsocketUpgrader = websocket.Upgrader{
	CheckOrigin:     checkOrigin,
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var (
	ErrEventNotSupported = errors.New("this event type is not supported")
)

type Manager struct {
	Clients      ClientList
	sync.RWMutex // to protect clients activity with mutex
	handlers     map[string]EventHandler
}

// NewManager creates new Manager
func NewManager() *Manager {
	m := &Manager{
		Clients:  make(ClientList),
		handlers: make(map[string]EventHandler),
	}
	m.setupEventHandlers()
	return m
}

// setupEventHandlers adds Event handlers to handlers Map
func (m *Manager) setupEventHandlers() {

	// may add different events e.g. for sending posts, comments, notifications in future
	m.handlers[EventMessage] = m.handleMessages
}

// handleMessages takes care of sent messages, save later to DB here
func (m *Manager) handleMessages(e Event, c *Client) error {

	var message structs.Message
	fmt.Printf("Handling %v event\n", e.Type)

	err := json.Unmarshal(e.Payload, &message)
	if err != nil {
		return fmt.Errorf("error unmarshalling the payload: %w", err)
	}
	fmt.Printf("New message: %v Created at: %v\n", message.Content, message.Created)

	return nil
}

// routeEvent routing Events to appropriate handler and validate sessions
func (m *Manager) routeEvent(event Event, c *Client) error {

	// Check if Handler is present in Map
	if handler, ok := m.handlers[event.Type]; ok {
		// Execute the handler and return any err
		if err := handler(event, c); err != nil {
			return err
		}
		return nil
	} else {
		return ErrEventNotSupported
	}
}

// Serve_WS upgrading regular http connection into websocket
func (m *Manager) Serve_WS(w http.ResponseWriter, r *http.Request) {

	// Begin by upgrading the HTTP request
	conn, err := WebsocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	// Create New Client
	client := NewClient(conn, m)
	// Add the newly created client to the manager
	m.addClient(client)

	go client.readMessages()
	go client.writeMessages()

}

// addClient is concurrently adding client to manager (w/mutex)
func (m *Manager) addClient(client *Client) {
	m.Lock()
	defer m.Unlock()
	m.Clients[client] = true
}

// removeClient concurrently safely (w/mutex) removing client
func (m *Manager) removeClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	// Check if Client exists, then delete it
	if _, ok := m.Clients[client]; ok {
		// close connection
		client.connection.Close()
		// remove
		//fmt.Println(client.nickname, "WS connection closed")
		delete(m.Clients, client)
	}
}

// checkOrigin will check origin and return true if it's allowed
func checkOrigin(r *http.Request) bool {
	// Grab the request origin
	origin := r.Header.Get("Origin")

	switch origin {
	case "http://localhost:8080", "http://localhost:3000":
		return true
	default:
		return false
	}
}
