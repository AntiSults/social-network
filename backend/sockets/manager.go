package sockets

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"social-network/db/sqlite"
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

var managerInstance *Manager
var once sync.Once

// GetManager returns the singleton instance of the Manager
func GetManager() *Manager {
	once.Do(func() {
		managerInstance = NewManager()
	})
	return managerInstance
}

// setupEventHandlers adds Event handlers to handlers Map
func (m *Manager) setupEventHandlers() {

	// may add different events e.g. for sending posts, comments, notifications in future
	m.handlers[EventMessage] = m.handleMessages
	m.handlers[EventGroupMessage] = m.handleMessages
	m.handlers[EventUpload] = m.handleUpload
	m.handlers[EventGroupUpload] = m.handleUpload
	m.handlers[EventNotify] = m.HandleNotify
}

func (m *Manager) HandleNotify(e Event, c *Client) error {
	var follower structs.User // assuming User struct is defined elsewhere
	err := json.Unmarshal(e.Payload, &follower)
	if err != nil {
		log.Printf("Error unmarshaling follower info: %v", err)
		return err
	}
	// Prepare notification to be sent
	notifyEvent := newEvent(EventNotify, e.Payload, e.SessionToken)

	// Send notification to the user (follower)
	if client, ok := m.ClientsByUserID[c.clientId]; ok {
		client.egress <- *notifyEvent

	} else {
		log.Printf("User with ID %d not connected", c.clientId)
		return fmt.Errorf("user not connected")
	}

	return nil
}
func (m *Manager) handleUpload(e Event, c *Client) error {
	if e.Type == "initial_group_upload" {
		var req struct {
			GroupID string `json:"groupId"`
		}
		err := json.Unmarshal(e.Payload, &req)
		if err != nil {
			return fmt.Errorf("error unmarshalling the payload (GroupID): %w", err)
		}
		groupID, err := strconv.Atoi(req.GroupID)
		if err != nil || groupID <= 0 {
			return fmt.Errorf("invalid groupID")
		}
		token := e.SessionToken
		if token == "" {
			return fmt.Errorf("session token is missing")
		}
		// Attempt to get the userID from in-memory session store, then DB
		userID, err := middleware.GetUserId(token)
		if err != nil {
			return fmt.Errorf("error getting ID from session token: %w", err)
		}
		// Fetch group messages for current User & GroupID
		messages, err := sqlite.Db.FetchGroupMessages(groupID, userID)
		if err != nil {
			return fmt.Errorf("error fetching messages for user ID %d: %w", userID, err)
		}
		fmt.Printf("GroupID %d,\nUserID %d, \nGroup messages %v\n", groupID, userID, messages)
		// Fetch the user IDs of the members in the group
		groupUserIDs, err := sqlite.Db.GetGroupUsers(userID, groupID)
		if err != nil {
			return fmt.Errorf("error querying group users: %w", err)
		}
		//add current user for the groupIDs
		groupUserIDs = append(groupUserIDs, userID)
		// Fetch full user information from the user IDs
		usersInfo, err := sqlite.Db.GetUsersByIDs(groupUserIDs)
		if err != nil {
			return fmt.Errorf("error querying usersInfo: %w", err)
		}
		// Create the response payload with messages and full user information
		common := structs.ChatMessage{
			Message: messages,
			User:    usersInfo,
		}
		// Marshal the response to JSON
		dataJSON, err := json.Marshal(&common)
		if err != nil {
			log.Println("error marshaling messages: ", err)
			return err
		}
		// Create the response event for group upload
		updateEvent := newEvent("initial_group_upload_response", dataJSON, token)
		// Send the response to the correct client
		if client, ok := m.ClientsByUserID[userID]; ok {
			client.egress <- *updateEvent
		}
		return nil
	}
	// Regular initial upload (for non-group chat)
	if e.Type == "initial_upload" {
		token := e.SessionToken
		if token == "" {
			return fmt.Errorf("session token is missing")
		}
		// Attempt to get the userID from in-memory session store, then DB
		userID, err := middleware.GetUserId(token)
		if err != nil {
			return fmt.Errorf("error getting ID from session token: %w", err)
		}
		// Fetch messages for the current user
		messages, err := sqlite.Db.FetchMessages(userID)
		if err != nil {
			return fmt.Errorf("error fetching messages for user ID %d: %w", userID, err)
		}
		// Fetch the user's followers and their info
		followerSlice, err := sqlite.Db.GetFollowersSlice(userID)
		if err != nil {
			return fmt.Errorf("error querying followers slice data: %w", err)
		}
		// Including the current user
		followerSlice = append(followerSlice, userID)
		// Fetch user info from the DB
		usersInfo, err := sqlite.Db.GetUsersByIDs(followerSlice)
		if err != nil {
			return fmt.Errorf("error querying usersInfo: %w", err)
		}
		// Create the response payload with messages and users
		common := structs.ChatMessage{
			Message: messages,
			User:    usersInfo,
		}
		// Marshal the response to JSON
		dataJSON, err := json.Marshal(&common)
		if err != nil {
			log.Println("error marshaling messages: ", err)
			return err
		}
		// Create the response event for regular chat upload
		updateEvent := newEvent("initial_upload_response", dataJSON, token)
		// Send the response to the correct client
		if client, ok := m.ClientsByUserID[userID]; ok {
			client.egress <- *updateEvent
		}
		return nil
	}
	// Handle unexpected event types
	return fmt.Errorf("unexpected event type: %s", e.Type)
}

// handleMessages takes care of sent messages, save later to DB here
func (m *Manager) handleMessages(e Event, c *Client) error {
	if e.Type == "group_chat_message" {
		var common structs.ChatMessage
		fmt.Printf("Handling %v event\n", e.Type)

		err := json.Unmarshal(e.Payload, &common)
		if err != nil {
			return fmt.Errorf("error unmarshalling the payload: %w", err)
		}
		_, err = sqlite.Db.SaveGroupMessage(&common.Message[0])
		if err != nil {
			log.Println("error saving PM into db: ", err)
		}
		// finding user to send message to
		updateEvent := newEvent("chat_message", e.Payload, e.SessionToken)

		if recipientClient, ok := m.ClientsByUserID[common.Message[0].RecipientID]; ok {
			recipientClient.egress <- *updateEvent
		}
		return nil
	} else if e.Type == "chat_message" {
		var common structs.ChatMessage
		fmt.Printf("Handling %v event\n", e.Type)
		err := json.Unmarshal(e.Payload, &common)
		if err != nil {
			return fmt.Errorf("error unmarshalling the payload: %w", err)
		}

		_, err = sqlite.Db.SaveMessage(&common.Message[0])

		if err != nil {
			log.Println("error saving PM into db: ", err)
		}
		// finding user to send message to
		updateEvent := newEvent("chat_message", e.Payload, e.SessionToken)

		if recipientClient, ok := m.ClientsByUserID[common.Message[0].RecipientID]; ok {
			recipientClient.egress <- *updateEvent
		}
		return nil
	}
	// Handle unexpected event types
	return fmt.Errorf("unexpected event type: %s", e.Type)
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
	userID, err := middleware.GetUserId(cookie.Value)
	if err != nil {
		middleware.SendErrorResponse(w, "Error getting user ID: "+err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Printf("Client %d is connected\n", userID)
	// Begin by upgrading the HTTP request
	conn, err := WebsocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("Client %d is connecting", userID)
	// Create New Client
	client := NewClient(conn, m)
	client.clientId = userID
	// Add the newly created client to the manager
	m.addClient(client)
	m.ClientsByUserID[userID] = client

	log.Printf("Client %d added. Current Clients: %+v", userID, m.Clients)
	log.Printf("Clients By UserID: %+v", m.ClientsByUserID)

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
