package security

import (
	"fmt"
	"net/http"
	"social-network/db/sqlite"
	"sync"
	"time"

	"github.com/satori/go.uuid"
)

const sessionLength int = 1800 // seconds

var (
	dbSessions  = map[string]Session{}
	sessionLock sync.Mutex
)

type Session struct {
	UserName       string
	UserID         int
	SessionToken   string
	LastActivity   time.Time
	ExpirationTime time.Time
}

func CleanSessions() {
	sessionLock.Lock()
	defer sessionLock.Unlock()

	for k, v := range dbSessions {
		if time.Since(v.LastActivity) > (time.Second * time.Duration(sessionLength)) {
			delete(dbSessions, k)
		}
	}
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
	sessionLock.Lock()
	defer sessionLock.Unlock()

	// Remove any existing session for the same user
	for token, session := range dbSessions {
		if session.UserID == userID {
			delete(dbSessions, token)
		}
	}
	dbSessions[token] = session
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
	sessionLock.Lock()
	defer sessionLock.Unlock()

	session, exists := dbSessions[sessionToken]
	if !exists {
		return false
	}

	// Check if the session is expired
	if session.ExpirationTime.Before(time.Now()) {
		delete(dbSessions, sessionToken)
		return false
	}

	// Update last activity
	session.LastActivity = time.Now()
	dbSessions[sessionToken] = session
	return true
}

// StartSessionCleaner periodically clean sessions
func StartSessionCleaner() {
	ticker := time.NewTicker(time.Duration(sessionLength) * time.Second)
	go func() {
		for range ticker.C {
			CleanSessions()
		}
	}()
}
