package sockets

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"social-network/db/sqlite"
	"social-network/handlers"
	"social-network/middleware"
	"social-network/structs"
	"strconv"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
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
	Clients         ClientList
	sync.RWMutex    // to protect clients activity with mutex
	ClientsByUserID map[int]*Client
	handlers        map[string]EventHandler
}

// NewManager creates new Manager
func NewManager() *Manager {
	m := &Manager{
		Clients:         make(ClientList),
		ClientsByUserID: make(map[int]*Client),
		handlers:        make(map[string]EventHandler),
	}
	m.setupEventHandlers()
	return m
}

// setupEventHandlers adds Event handlers to handlers Map
func (m *Manager) setupEventHandlers() {

	// may add different events e.g. for sending posts, comments, notifications in future
	m.handlers[EventMessage] = m.handleMessages
	m.handlers[EventUpload] = m.handleUpload
}

func (m *Manager) handleUpload(e Event, c *Client) error {
	fmt.Println("Request Upload type", e.Type)
	if e.Type != "initial_upload" {
		return fmt.Errorf("unexpected event type: %s", e.Type)
	}
	token := e.SessionToken
	if token == "" {
		return fmt.Errorf("session token is missing")
	}
	// Attempt to get the userID from in-memory session store, then DB
	userID, err := handlers.GetUserId(token)
	if err != nil {
		return fmt.Errorf("error getting ID from session token: %w", err)
	}
	// getting slice of followers
	followerSlice, err := sqlite.Db.GetFollowersSlice(userID)
	if err != nil {
		return fmt.Errorf("error querying followers slice data: %w", err)
	}
	//including current user
	followerSlice = append(followerSlice, userID)
	//getting users from db
	usersInfo, err := sqlite.Db.GetUsersByIDs(followerSlice)
	if err != nil {
		return fmt.Errorf("error querying usersInfo: %w", err)
	}
	// Fetch messages for current User
	messages, err := sqlite.Db.FetchMessages(userID)
	if err != nil {
		return fmt.Errorf("error fetching messages for user ID %d: %w", userID, err)
	}
	common := structs.ChatMessage{
		Message: messages,
		User:    usersInfo,
		// Group:   groupInfo,
	}
	// Marshal messages to JSON
	dataJSON, err := json.Marshal(&common)
	if err != nil {
		log.Println("error marshaling messages: ", err)
		return err
	}
	// Create the response event
	updateEvent := newEvent("initial_upload_response", dataJSON, token)

	// Send the response to the correct client
	if client, ok := m.ClientsByUserID[userID]; ok {
		client.egress <- *updateEvent
	}
	return nil
}

// handleMessages takes care of sent messages, save later to DB here
func (m *Manager) handleMessages(e Event, c *Client) error {

	var common structs.ChatMessage
	fmt.Printf("Handling %v event\n", string(e.Type))
	err := json.Unmarshal(e.Payload, &common)
	if err != nil {
		return fmt.Errorf("error unmarshalling the payload: %w", err)
	}

	_, err = sqlite.Db.SaveMessage(&common.Message[0])

	if err != nil {
		log.Println("error saving PM into db: ", err)
	}
	// finding user to send message to
	updateEvent := newEvent("chat_message", e.Payload, "")

	if recipientClient, ok := m.ClientsByUserID[common.Message[0].RecipientID]; ok {
		recipientClient.egress <- *updateEvent
	}
	return nil
}

// routeEvent routing Events to appropriate handler
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

	cookie, err := r.Cookie("session_token")
	if err != nil {
		middleware.SendErrorResponse(w, "Error getting token: "+err.Error(), http.StatusBadRequest)
		return
	}
	userID, err := handlers.GetUserId(cookie.Value)
	if err != nil {
		middleware.SendErrorResponse(w, "Error getting user ID: "+err.Error(), http.StatusBadRequest)
		return
	}
	// Begin by upgrading the HTTP request
	conn, err := WebsocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	// Create New Client
	client := NewClient(conn, m)
	client.clientId = userID
	// Add the newly created client to the manager
	m.addClient(client)
	m.ClientsByUserID[userID] = client

	go client.readMessages()
	go client.writeMessages()

}

// addClient is concurrently adding client to manager (w/mutex)
func (m *Manager) addClient(client *Client) {
	m.Lock()
	defer m.Unlock()
	m.Clients[client] = true
	m.ClientsByUserID[client.clientId] = client
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
		delete(m.Clients, client)
		delete(m.ClientsByUserID, client.clientId) // Remove from ClientsByUserID
	}
}

// checkOrigin will check origin and return true if it's allowed
func checkOrigin(r *http.Request) bool {
	// Grab the request origin
	origin := r.Header.Get("Origin")
	if origin == "" {
		return false
	}
	// Parse the origin to extract the host and port
	u, err := url.Parse(origin)
	if err != nil {
		return false
	}
	hostParts := strings.Split(u.Host, ":")
	if len(hostParts) != 2 {
		return false
	}
	port, err := strconv.Atoi(hostParts[1])
	if err != nil {
		return false
	}
	// Allow localhost with port range 3000-3010 or specific other origins
	if (hostParts[0] == "localhost" && port >= 3000 && port <= 3010) || origin == "http://localhost:8080" {
		return true
	}
	return false
}
