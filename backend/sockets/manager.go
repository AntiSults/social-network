package sockets

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"social-network/db/sqlite"
	"social-network/handlers"
	"social-network/middleware"
	"social-network/structs"
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
	m.handlers[EventUpload] = m.handleUpload
}

func (m *Manager) handleUpload(e Event, c *Client) error {

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
	//Getting user either from map or db
	user, err := handlers.GetUser(userID)
	if err != nil {
		return fmt.Errorf("error querying user data: %w", err)
	}
	//hardcoding followers, as not yet implimented
	user.FollowingUserIDs = []int{2, 3, 5}
	user.GotFollowedUserIDs = []int{3, 4}

	//combining followers and followed into single slice
	usersID := combineUnique(user.FollowingUserIDs, user.GotFollowedUserIDs)
	//including current user
	usersID = append(usersID, userID)

	//getting users from db
	usersInfo, err := sqlite.Db.GetUsersByIDs(usersID)
	if err != nil {
		return fmt.Errorf("error querying usersInfo: %w", err)
	}
	// This is appending followers and groups into usersInfo, which in turn is send
	//with messages to frontend as initial upload response. But it is not used there for now
	for i := range usersInfo {
		if usersInfo[i].ID == userID {
			usersInfo[i].FollowingUserIDs = user.FollowingUserIDs
			usersInfo[i].GotFollowedUserIDs = user.GotFollowedUserIDs
		}
	}
	// Fetch messages
	messages, err := sqlite.Db.FetchMessages(userID)
	if err != nil {
		return fmt.Errorf("error fetching messages for user ID %d: %w", userID, err)
	}
	common := structs.ChatMessage{
		Message: messages,
		User:    usersInfo,
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
	for client := range m.Clients {
		if userID == c.clientId {
			client.egress <- *updateEvent
			break // Exit the loop once the correct client is found
		}
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
	fmt.Println("New message:", &common.Message[0])

	_, err = sqlite.Db.SaveMessage(&common.Message[0])

	if err != nil {
		log.Println("error saving PM into db: ", err)
	}
	// finding user or users to send message to
	updateEvent := newEvent("message_received", e.Payload, "")
	for client := range m.Clients {
		for recipient := range common.Message[0].RecipientID {
			if recipient == c.clientId {
				client.egress <- *updateEvent
			}
		}
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
	// Add the newly created client to the manager
	m.addClient(client)
	client.clientId = userID

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

func combineUnique(slice1, slice2 []int) []int {
	uniqueMap := make(map[int]bool)
	var result []int
	for _, v := range slice1 {
		if !uniqueMap[v] {
			uniqueMap[v] = true
			result = append(result, v)
		}
	}
	for _, v := range slice2 {
		if !uniqueMap[v] {
			uniqueMap[v] = true
			result = append(result, v)
		}
	}
	return result
}
