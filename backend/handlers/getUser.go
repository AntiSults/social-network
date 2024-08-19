package handlers

import (
	"encoding/json"
	"net/http"

	"social-network/db/sqlite"
	"social-network/middleware"
)

func GetUserData(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {

		cookie, err := r.Cookie("session_token")
		if err != nil {
			middleware.SendErrorResponse(w, "Error getting token"+err.Error(), http.StatusBadRequest)
			return
		}

		userID, err := sqlite.Db.GetUserIdFromToken(cookie.Value)
		if err != nil {
			middleware.SendErrorResponse(w, "Error getting ID from session token", http.StatusInternalServerError)
			return
		}

		user, err := sqlite.Db.GetUser(userID)

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
