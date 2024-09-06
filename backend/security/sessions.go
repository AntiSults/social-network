package security

import (
	"fmt"
	"github.com/satori/go.uuid"
	"net/http"
	"social-network/db/sqlite"
	"sync"
	"time"
)

const sessionLength int = 1800 // seconds

var (
	DbSessions  = map[string]Session{}
	SessionLock sync.RWMutex
)

type Session struct {
	UserName       string
	UserID         int
	SessionToken   string
	LastActivity   time.Time
	ExpirationTime time.Time
}

func RemoveSession(token string) {
	SessionLock.Lock()
	defer SessionLock.Unlock()
	delete(DbSessions, token)
}

func NewSession(userName string, userID int, w http.ResponseWriter) {
	token := uuid.NewV4().String()
	session := Session{
		UserName:       userName,
		UserID:         userID,
		SessionToken:   token,
		LastActivity:   time.Now(),
		ExpirationTime: time.Now().Add(time.Second * time.Duration(sessionLength)),
	}
	SessionLock.Lock()
	defer SessionLock.Unlock()

	// Remove any existing session for the same user
	for token, session := range DbSessions {
		if session.UserID == userID {
			delete(DbSessions, token)
		}
	}
	DbSessions[token] = session
	session.setCookie(w)
	err := sqlite.Db.SaveSession(userID, token, session.ExpirationTime)

	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error inserting a session", http.StatusInternalServerError)
		return
	}
}

func (s *Session) setCookie(w http.ResponseWriter) {
	cookie := http.Cookie{
		Name:    s.UserName,
		Path:    "/",
		Value:   s.SessionToken,
		Expires: s.ExpirationTime,
		MaxAge:  sessionLength,
	}
	http.SetCookie(w, &cookie)
}
func ValidateSession(sessionToken string) bool {
	SessionLock.Lock()
	defer SessionLock.Unlock()

	session, exists := DbSessions[sessionToken]
	if !exists {
		return false
	}

	// Check if the session is expired
	if session.ExpirationTime.Before(time.Now()) {
		delete(DbSessions, sessionToken)
		return false
	}

	// Update last activity
	session.LastActivity = time.Now()
	DbSessions[sessionToken] = session
	return true
}

// StartSessionCleaner periodically clean sessions
func (s *Session) StartSessionCleaner() {
	ticker := time.NewTicker(time.Duration(sessionLength) * time.Second)
	go func() {
		for range ticker.C {
			RemoveSession(s.SessionToken)
		}
	}()
}
