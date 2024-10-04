package handlers

import (
	"encoding/json"
	"net/http"
	"social-network/middleware"
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
	userID, err := middleware.GetUserId(cookie.Value)
	if err != nil {
		middleware.SendErrorResponse(w, "Error getting ID from session token", http.StatusInternalServerError)
		return
	}
	user, err := middleware.GetUser(userID)

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
