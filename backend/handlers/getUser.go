package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"social-network/db/sqlite"
	"social-network/middleware"
	"social-network/security"
	"social-network/structs"
)

func GetUserData(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		middleware.SendErrorResponse(w, "Method not allowed!", http.StatusMethodNotAllowed)
	}
	cookie, err := r.Cookie("session_token")
	if err != nil {
		middleware.SendErrorResponse(w, "Error getting token"+err.Error(), http.StatusBadRequest)
		return
	}

	userID, err := GetUserId(cookie.Value)
	if err != nil {
		middleware.SendErrorResponse(w, "Error getting ID from session token", http.StatusInternalServerError)
		return
	}

	user, err := GetUser(userID)

	if err != nil {
		middleware.SendErrorResponse(w, "Error querying user data to struct"+err.Error(), http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(user)
	if err != nil {
		middleware.SendErrorResponse(w, "Error marshalling user data to JSON", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)

}

// GetUserId is getting user ID with token, either from sessions map or from DB
func GetUserId(token string) (int, error) {
	userID := 0
	var err error

	// Acquire a read lock before accessing the shared map
	security.SessionLock.RLock()
	session, ok := security.DbSessions[token]
	security.SessionLock.RUnlock()

	if ok {
		userID = session.UserID
	} else {
		// Fall back to database lookup if not found in in-memory store
		userID, err = sqlite.Db.GetUserIdFromToken(token)
		if err != nil {
			return -1, fmt.Errorf("error getting ID from session token: %w", err)
		}
	}
	return userID, err
}

// GetUser is getting user with user ID, either from User map or from DB
func GetUser(id int) (*structs.User, error) {
	var (
		user *structs.User
		err  error
	)
	UserMapLock.RLock()
	u, ok := UserMap[id]
	UserMapLock.RUnlock()

	if ok {
		user = &u
	} else {
		// Fall back to database lookup if not found in in-memory store
		user, err = sqlite.Db.GetUser(id)
		if err != nil {
			return nil, fmt.Errorf("error querying user data to struct:  %w", err)
		}
	}
	return user, nil
}
