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
	Clients ClientList
	sync.RWMutex
	ClientsByUserID map[int]map[string]*Client
	handlers        map[string]EventHandler
}

// NewManager creates new Manager
func NewManager() *Manager {
	m := &Manager{
		Clients:         make(ClientList),
		ClientsByUserID: make(map[int]map[string]*Client),
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
}

func (m *Manager) HandleNotify(e Event, c *Client) error {

	notifyEvent := newEvent(e.Type, e.Payload, e.SessionToken)

	if e.Type == EventNewGroupEvent {
		var req struct {
			GroupEvent []structs.Event
			GroupName  string
		}
		err := json.Unmarshal(e.Payload, &req)
		if err != nil {
			return fmt.Errorf("error unmarshalling the payload fordata struct: %w", err)
		}
		userIDs, err := sqlite.Db.GetGroupUserIDs(req.GroupEvent[0].GroupID)
		if err != nil {
			return fmt.Errorf("error getting slice of userID from GroupUsers: %w", err)
		}
		for _, userID := range userIDs {
			if client, ok := m.ClientsByUserID[userID]["notify"]; ok {
				client.egress <- *notifyEvent
			}
		}
		return nil
	}
	if client, ok := m.ClientsByUserID[c.clientId]["notify"]; ok {
		client.egress <- *notifyEvent
	} else {
		log.Printf("User with ID %d not connected", c.clientId)
		return fmt.Errorf("user not connected")
	}
	return nil
}
func (m *Manager) handleUpload(e Event, c *Client) error {
	if e.Type == EventGroupUpload {
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
		if client, ok := m.ClientsByUserID[userID]["chat"]; ok {
			client.egress <- *updateEvent
		}
		return nil
	}
	// Regular initial upload (for non-group chat)
	if e.Type == EventUpload {
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
		if client, ok := m.ClientsByUserID[userID]["chat"]; ok {
			client.egress <- *updateEvent
		}
		return nil
	}
	// Handle unexpected event types
	return fmt.Errorf("unexpected event type: %s", e.Type)
}

// handleMessages takes care of sent messages, save later to DB here
func (m *Manager) handleMessages(e Event, c *Client) error {
	if e.Type == EventGroupMessage {
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
		updateEvent := newEvent(EventGroupMessage, e.Payload, e.SessionToken)

		if recipientClient, ok := m.ClientsByUserID[common.Message[0].RecipientID]["chat"]; ok {
			recipientClient.egress <- *updateEvent
		}
		return nil
	} else if e.Type == EventMessage {
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
		updateEvent := newEvent(EventMessage, e.Payload, e.SessionToken)

		if recipientClient, ok := m.ClientsByUserID[common.Message[0].RecipientID]["chat"]; ok {
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
	// Determine connection type
	connType := r.URL.Path
	if connType == "/ws" {
		connType = "chat"
	} else if connType == "/notify" {
		connType = "notify"
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
	client.connType = connType
	// Add the newly created client to the manager
	m.addClient(client, connType, userID)

	go client.readMessages()
	go client.writeMessages()
}

// addClient is concurrently adding client to manager (w/mutex)
func (m *Manager) addClient(client *Client, connType string, userID int) {
	m.Lock()
	defer m.Unlock()

	// Check if the user already has a connection map; if not, create it
	if _, ok := m.ClientsByUserID[userID]; !ok {
		m.ClientsByUserID[userID] = make(map[string]*Client)
	}

	// Add the new client under the correct connection type
	m.ClientsByUserID[userID][connType] = client

	// Add to global client list
	m.Clients[client] = true
}

// removeClient concurrently safely (w/mutex) removing client
// removeClient safely removes the client from the manager
func (m *Manager) removeClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	userID := client.clientId
	connType := client.connType

	// Check if the user has a connection of this type and remove it
	if _, ok := m.ClientsByUserID[userID][connType]; ok {
		client.connection.Close()
		delete(m.ClientsByUserID[userID], connType)
		delete(m.Clients, client)

		// Optionally, you can remove the user's entry entirely if no connections remain
		if len(m.ClientsByUserID[userID]) == 0 {
			delete(m.ClientsByUserID, userID)
		}
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
