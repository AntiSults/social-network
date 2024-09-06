package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"social-network/db/sqlite"
	"social-network/middleware"
	"social-network/security"
)

func GetUserData(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {

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

		user, err := sqlite.Db.GetUser(userID)
		fmt.Println(user)

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
	} else {
		middleware.SendErrorResponse(w, "Method not allowed!", http.StatusMethodNotAllowed)
	}
}

func GetUserId(token string) (int, error) {
	userID := 0
	var err error
	if session, ok := security.DbSessions[token]; ok {
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
